package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Wenfuuu/tpa-web/pkg/messaging"
	"github.com/Wenfuuu/tpa-web/service-email/utils"
	"github.com/rabbitmq/amqp091-go"
	"gopkg.in/gomail.v2"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func setupVerificationEmailQueue(ch *amqp091.Channel, exchangeName string) amqp091.Queue {
	queueName := "verification_email_queue"
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare verification queue")

	routingKey := "verification_email"
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind verification queue")

	return q
}

func setupEmailQueue(ch *amqp091.Channel, exchangeName string) amqp091.Queue {
	queueName := "email_queue"
	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare normal email queue")

	routingKey := "email"
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind normal email queue")

	return q
}

func setupDeadLetterQueue(ch *amqp091.Channel, exchangeName string) amqp091.Queue {
	deadLetterQueueName := "failed_emails_queue"
	q, err := ch.QueueDeclare(
		deadLetterQueueName, // name
		true,                // durable
		false,               // delete when unused
		false,               // exclusive
		false,               // no-wait
		nil,                 // arguments
	)
	failOnError(err, "Failed to declare dead letter queue")

	routingKey := "failed_email"
	err = ch.QueueBind(
		q.Name,       // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind dead letter queue")

	return q
}

type RetryMessage struct {
	OriginalMessage []byte `json:"original_message"`
	RetryCount      int    `json:"retry_count"`
	MessageType     string `json:"message_type"` // "verification" or "normal"
}

func main() {
	log.Println("Starting Email Service...")

	ch, err := utils.GetRabbitMQChannel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchangeName := "email_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, "direct", true, false, false, false, nil,
	)
	failOnError(err, "Failed to declare exchange")

	verifQueue := setupVerificationEmailQueue(ch, exchangeName)
	emailQueue := setupEmailQueue(ch, exchangeName)
	deadLetterQueue := setupDeadLetterQueue(ch, exchangeName)

	retryQueueName := "email_retry_queue"
	retryQueue, err := ch.QueueDeclare(
		retryQueueName, // name
		true,           // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare retry queue")

	err = ch.QueueBind(
		retryQueue.Name, // queue name
		"email_retry",   // routing key
		exchangeName,    // exchange
		false,
		nil,
	)
	failOnError(err, "Failed to bind retry queue")

	err = ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	verificationMsgs, err := ch.Consume(
		verifQueue.Name, // queue
		"",              // consumer tag (empty means auto-generate)
		false,           // auto-ack (set to false for manual ack)
		false,           // exclusive
		false,           // no-local (not supported on RabbitMQ)
		false,           // no-wait
		nil,             // args
	)
	failOnError(err, "Failed to register verification email consumer")

	normalMsgs, err := ch.Consume(
		emailQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register normal email consumer")

	retryMsgs, err := ch.Consume(
		retryQueue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register retry consumer")
	
	deadLetterMsgs, err := ch.Consume(
		deadLetterQueue.Name,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register dead letter queue consumer")

	const maxRetries = 3

	go func() {
		for d := range verificationMsgs {
			log.Printf("Received a verification message: %s", d.Body)
			var payload messaging.VerificationEmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("ERROR: Failed to unmarshal message: %v", err)
				// send to dead letter queue
				publishToDeadLetterQueue(ch, exchangeName, d.Body, "verification", "Format error")
				d.Ack(false) 
				continue
			}

			if err := sendVerificationEmail(payload.Email, payload.Code); err != nil {
				log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
				
				retryMsg := RetryMessage{
					OriginalMessage: d.Body,
					RetryCount:      0,
					MessageType:     "verification",
				}
				
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false)
			} else {
				log.Printf("Successfully sent verification email to %s", payload.Email)
				d.Ack(false)
			}
		}
	}()

	go func() {
		for d := range normalMsgs {
			log.Printf("Received an email message: %s", d.Body)
			var payload messaging.EmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("ERROR: Failed to unmarshal message: %v", err)
				// send to dead letter queue
				publishToDeadLetterQueue(ch, exchangeName, d.Body, "normal", "Format error")
				d.Ack(false)
				continue
			}

			if err := sendEmail(payload.Email, payload.Word); err != nil {
				log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
				
				retryMsg := RetryMessage{
					OriginalMessage: d.Body,
					RetryCount:      0,
					MessageType:     "normal",
				}
				
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false)
			} else {
				log.Printf("Successfully sent email to %s", payload.Email)
				d.Ack(false)
			}
		}
	}()

	go func() {
		for d := range retryMsgs {
			log.Printf("Processing retry message: %s", d.Body)
			var retryMsg RetryMessage
			if err := json.Unmarshal(d.Body, &retryMsg); err != nil {
				log.Printf("ERROR: Failed to unmarshal retry message: %v", err)
				publishToDeadLetterQueue(ch, exchangeName, d.Body, "retry", "Format error")
				d.Ack(false)
				continue
			}

			if retryMsg.RetryCount >= maxRetries {
				log.Printf("Max retries reached for message, sending to dead letter queue")
				publishToDeadLetterQueue(ch, exchangeName, retryMsg.OriginalMessage, retryMsg.MessageType, "Max retries exceeded")
				d.Ack(false)
				continue
			}

			retryMsg.RetryCount++

			var err error
			if retryMsg.MessageType == "verification" {
				var payload messaging.VerificationEmailPayload
				if jsonErr := json.Unmarshal(retryMsg.OriginalMessage, &payload); jsonErr != nil {
					log.Printf("ERROR: Failed to unmarshal original verification message: %v", jsonErr)
					publishToDeadLetterQueue(ch, exchangeName, retryMsg.OriginalMessage, retryMsg.MessageType, "Format error")
					d.Ack(false)
					continue
				}
				err = sendVerificationEmail(payload.Email, payload.Code)
			} else { // normal email
				var payload messaging.EmailPayload
				if jsonErr := json.Unmarshal(retryMsg.OriginalMessage, &payload); jsonErr != nil {
					log.Printf("ERROR: Failed to unmarshal original email message: %v", jsonErr)
					publishToDeadLetterQueue(ch, exchangeName, retryMsg.OriginalMessage, retryMsg.MessageType, "Format error")
					d.Ack(false)
					continue
				}
				err = sendEmail(payload.Email, payload.Word)
			}

			if err != nil {
				log.Printf("Retry %d failed, will retry again: %v", retryMsg.RetryCount, err)
				
				delay := time.Duration(retryMsg.RetryCount * retryMsg.RetryCount) * time.Second
				time.Sleep(delay)
				
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false)
			} else {
				log.Printf("Retry %d successful", retryMsg.RetryCount)
				d.Ack(false)
			}
		}
	}()

	go func() {
		for d := range deadLetterMsgs {
			log.Printf("Processing dead letter message: %s", d.Body)
			
			var deadLetterMsg struct {
				OriginalMessage []byte    `json:"original_message"`
				MessageType     string    `json:"message_type"`
				FailureReason   string    `json:"failure_reason"`
				FailedAt        time.Time `json:"failed_at"`
			}
			
			if err := json.Unmarshal(d.Body, &deadLetterMsg); err != nil {
				log.Printf("ERROR: Failed to unmarshal dead letter message: %v", err)
				d.Ack(false)
				continue
			}
			
			log.Printf("DEAD LETTER: Type: %s, Reason: %s, Failed at: %s", 
				deadLetterMsg.MessageType, 
				deadLetterMsg.FailureReason,
				deadLetterMsg.FailedAt)
				
			
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("Shutting down Email Service...")
}

func publishToRetryQueue(ch *amqp091.Channel, exchangeName string, retryMsg RetryMessage) {
	retryMsgBytes, err := json.Marshal(retryMsg)
	if err != nil {
		log.Printf("ERROR: Failed to marshal retry message: %v", err)
		return
	}

	err = ch.Publish(
		exchangeName,  // exchange
		"email_retry", // routing key
		false,         // mandatory
		false,         // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         retryMsgBytes,
			DeliveryMode: amqp091.Persistent,
		})

	if err != nil {
		log.Printf("ERROR: Failed to publish to retry queue: %v", err)
	}
}

func publishToDeadLetterQueue(ch *amqp091.Channel, exchangeName string, originalMsg []byte, msgType string, reason string) {
	failedMsg := struct {
		OriginalMessage []byte    `json:"original_message"`
		MessageType     string    `json:"message_type"`
		FailureReason   string    `json:"failure_reason"`
		FailedAt        time.Time `json:"failed_at"`
	}{
		OriginalMessage: originalMsg,
		MessageType:     msgType,
		FailureReason:   reason,
		FailedAt:        time.Now(),
	}

	failedMsgBytes, err := json.Marshal(failedMsg)
	if err != nil {
		log.Printf("ERROR: Failed to marshal dead letter message: %v", err)
		return
	}

	err = ch.Publish(
		exchangeName,   // exchange
		"failed_email", // routing key
		false,          // mandatory
		false,          // immediate
		amqp091.Publishing{
			ContentType:  "application/json",
			Body:         failedMsgBytes,
			DeliveryMode: amqp091.Persistent,
		})

	if err != nil {
		log.Printf("ERROR: Failed to publish to dead letter queue: %v", err)
	} else {
		log.Printf("Published message to dead letter queue")
	}
}

func sendVerificationEmail(toEmail, code string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", utils.GetDotEnv("EMAIL_USER"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Your Verification Code")
	m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s. It expires in 5 minutes.", code))

	host := utils.GetDotEnv("SMTP_HOST")
	portStr := utils.GetDotEnv("SMTP_PORT")
	email := utils.GetDotEnv("EMAIL_USER")
	pass := utils.GetDotEnv("EMAIL_PASS")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	d := gomail.NewDialer(host, port, email, pass)

	err = d.DialAndSend(m)
	return err
}

func sendEmail(toEmail, word string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", utils.GetDotEnv("EMAIL_USER"))
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Notification")
	m.SetBody("text/plain", word)

	host := utils.GetDotEnv("SMTP_HOST")
	portStr := utils.GetDotEnv("SMTP_PORT")
	email := utils.GetDotEnv("EMAIL_USER")
	pass := utils.GetDotEnv("EMAIL_PASS")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP_PORT: %v", err)
	}

	d := gomail.NewDialer(host, port, email, pass)

	err = d.DialAndSend(m)
	return err
}