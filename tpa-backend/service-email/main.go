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

// setupVerificationEmailQueue declares and binds the verification email queue
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

// setupEmailQueue declares and binds the normal email queue
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

// setupDeadLetterQueue declares and sets up a dead letter queue for failed messages
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

// Message with retry information
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

	// Declare the same exchange as the publisher (idempotent)
	exchangeName := "email_exchange"
	err = ch.ExchangeDeclare(
		exchangeName, "direct", true, false, false, false, nil,
	)
	failOnError(err, "Failed to declare exchange")

	verifQueue := setupVerificationEmailQueue(ch, exchangeName)
	emailQueue := setupEmailQueue(ch, exchangeName)
	deadLetterQueue := setupDeadLetterQueue(ch, exchangeName)

	// Set up retry queue
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

	// Configure QoS: Process only one message at a time until acknowledged
	err = ch.Qos(
		1,     // prefetchCount
		0,     // prefetchSize
		false, // global
	)
	failOnError(err, "Failed to set QoS")

	// Start consuming messages
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
	
	// Consume from dead letter queue to process failed messages
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

	// Maximum number of retries before sending to dead letter queue
	const maxRetries = 3

	// Goroutine to handle verification emails
	go func() {
		for d := range verificationMsgs {
			log.Printf("Received a verification message: %s", d.Body)
			var payload messaging.VerificationEmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("ERROR: Failed to unmarshal message: %v", err)
				// Bad format - send to dead letter queue
				publishToDeadLetterQueue(ch, exchangeName, d.Body, "verification", "Format error")
				d.Ack(false) // Acknowledge to remove from queue
				continue
			}

			// Process the message (send email)
			if err := sendVerificationEmail(payload.Email, payload.Code); err != nil {
				log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
				
				// Publish to retry queue with retry count of 0
				retryMsg := RetryMessage{
					OriginalMessage: d.Body,
					RetryCount:      0,
					MessageType:     "verification",
				}
				
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false) // Acknowledge the original message
			} else {
				log.Printf("Successfully sent verification email to %s", payload.Email)
				d.Ack(false) // Acknowledge message successfully processed
			}
		}
	}()

	// Goroutine to handle normal emails
	go func() {
		for d := range normalMsgs {
			log.Printf("Received an email message: %s", d.Body)
			var payload messaging.EmailPayload
			if err := json.Unmarshal(d.Body, &payload); err != nil {
				log.Printf("ERROR: Failed to unmarshal message: %v", err)
				// Bad format - send to dead letter queue
				publishToDeadLetterQueue(ch, exchangeName, d.Body, "normal", "Format error")
				d.Ack(false) // Acknowledge to remove from queue
				continue
			}

			if err := sendEmail(payload.Email, payload.Word); err != nil {
				log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
				
				// Publish to retry queue with retry count of 0
				retryMsg := RetryMessage{
					OriginalMessage: d.Body,
					RetryCount:      0,
					MessageType:     "normal",
				}
				
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false) // Acknowledge the original message
			} else {
				log.Printf("Successfully sent email to %s", payload.Email)
				d.Ack(false) // Acknowledge message successfully processed
			}
		}
	}()

	// Goroutine to handle retries
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

			// Check retry count
			if retryMsg.RetryCount >= maxRetries {
				log.Printf("Max retries reached for message, sending to dead letter queue")
				publishToDeadLetterQueue(ch, exchangeName, retryMsg.OriginalMessage, retryMsg.MessageType, "Max retries exceeded")
				d.Ack(false)
				continue
			}

			// Increase retry count
			retryMsg.RetryCount++

			// Process based on message type
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
				
				// Implement exponential backoff with delay
				delay := time.Duration(retryMsg.RetryCount * retryMsg.RetryCount) * time.Second
				time.Sleep(delay)
				
				// Republish with incremented retry count
				publishToRetryQueue(ch, exchangeName, retryMsg)
				d.Ack(false)
			} else {
				log.Printf("Retry %d successful", retryMsg.RetryCount)
				d.Ack(false)
			}
		}
	}()

	// Goroutine to handle dead letter messages
	go func() {
		for d := range deadLetterMsgs {
			log.Printf("Processing dead letter message: %s", d.Body)
			
			// Parse the dead letter message
			var deadLetterMsg struct {
				OriginalMessage []byte    `json:"original_message"`
				MessageType     string    `json:"message_type"`
				FailureReason   string    `json:"failure_reason"`
				FailedAt        time.Time `json:"failed_at"`
			}
			
			if err := json.Unmarshal(d.Body, &deadLetterMsg); err != nil {
				log.Printf("ERROR: Failed to unmarshal dead letter message: %v", err)
				// This is a serious error, but we can't do much about it
				d.Ack(false)
				continue
			}
			
			// Log the failed message details for monitoring/alerting
			log.Printf("DEAD LETTER: Type: %s, Reason: %s, Failed at: %s", 
				deadLetterMsg.MessageType, 
				deadLetterMsg.FailureReason,
				deadLetterMsg.FailedAt)
				
			// You could implement additional logic here:
			// 1. Send alerts to admins
			// 2. Store in a database for later analysis
			// 3. Attempt special recovery procedures
			
			// For now, we're just acknowledging it
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c // Block until shutdown signal

	log.Println("Shutting down Email Service...")
}

// Publish a message to the retry queue
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

// Publish a message to the dead letter queue
func publishToDeadLetterQueue(ch *amqp091.Channel, exchangeName string, originalMsg []byte, msgType string, reason string) {
	// You might want to wrap the original message with some metadata
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
	// Use gomail to send the email
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