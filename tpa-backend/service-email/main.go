package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

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

// setupNormalEmailQueue declares and binds the normal email queue
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

    // Configure QoS: Process only one message at a time until acknowledged
    // Adjust prefetchCount based on your worker's capacity
    err = ch.Qos(
        1,     // prefetchCount
        0,     // prefetchSize
        false, // global
    )
    failOnError(err, "Failed to set QoS")

    // Start consuming messages
    verificationMsgs, err := ch.Consume(
        verifQueue.Name, // queue
        "",     // consumer tag (empty means auto-generate)
        false,  // auto-ack (set to false for manual ack)
        false,  // exclusive
        false,  // no-local (not supported on RabbitMQ)
        false,  // no-wait
        nil,    // args
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

    // Goroutine to verification emails
    go func() {
        for d := range verificationMsgs {
            log.Printf("Received a verification message: %s", d.Body)
            var payload messaging.VerificationEmailPayload
            if err := json.Unmarshal(d.Body, &payload); err != nil {
                log.Printf("ERROR: Failed to unmarshal message: %v", err)
                // Decide how to handle bad messages (e.g., move to dead-letter queue)
                d.Nack(false, false) // Nack message, don't requeue
                continue
            }

            // Process the message (send email)
            if err := sendVerificationEmail(payload.Email, payload.Code); err != nil {
                log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
                // Decide on requeue strategy based on error type (transient vs permanent)
                d.Nack(false, true) // Nack and requeue (use with caution, might cause infinite loop)
                                     // Consider delay or limited retries, or Nack(false, false) for permanent errors
            } else {
                log.Printf("Successfully sent verification email to %s", payload.Email)
                d.Ack(false) // Acknowledge message successfully processed
            }
        }
    }()

    // Gororutine to normal emails
    go func() {
        for d := range normalMsgs {
            log.Printf("Received a email message: %s", d.Body)
            var payload messaging.EmailPayload
            if err := json.Unmarshal(d.Body, &payload); err != nil {
                log.Printf("ERROR: Failed to unmarshal message: %v", err)
                d.Nack(false, false)
                continue
            }

            if err := sendEmail(payload.Email, payload.Word); err != nil {
                log.Printf("ERROR: Failed to send email to %s: %v", payload.Email, err)
                d.Nack(false, true)
            } else {
                log.Printf("Successfully sent email to %s", payload.Email)
                d.Ack(false) // Acknowledge message successfully processed
            }
        }
    }()

    log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

    // Graceful shutdown
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    <-c // Block until shutdown signal

    log.Println("Shutting down Email Service...")
}

func sendVerificationEmail(toEmail, code string) error {
    // Use gomail to send the email (same logic as before, but now here)
    m := gomail.NewMessage()
    m.SetHeader("From", utils.GetDotEnv("EMAIL_USER"))
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", "Your Verification Code")
    m.SetBody("text/plain", fmt.Sprintf("Your verification code is: %s. It expires in 5 minutes.", code))

	host := utils.GetDotEnv("SMTP_HOST") // e.g., "smtp.gmail.com"
    portStr := utils.GetDotEnv("SMTP_PORT") // e.g., "587"
    email := utils.GetDotEnv("EMAIL_USER")
    pass := utils.GetDotEnv("EMAIL_PASS")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

    d := gomail.NewDialer(host, port, email, pass)

    err = d.DialAndSend(m)
    return err // Return the error to the consumer loop
}

func sendEmail(toEmail, word string) error {
    m := gomail.NewMessage()
    m.SetHeader("From", utils.GetDotEnv("EMAIL_USER"))
    m.SetHeader("To", toEmail)
    m.SetHeader("Subject", "Notification")
    m.SetBody("text/plain", word)

    host := utils.GetDotEnv("SMTP_HOST") // e.g., "smtp.gmail.com"
    portStr := utils.GetDotEnv("SMTP_PORT") // e.g., "587"
    email := utils.GetDotEnv("EMAIL_USER")
    pass := utils.GetDotEnv("EMAIL_PASS")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid SMTP_PORT: %v", err)
	}

    d := gomail.NewDialer(host, port, email, pass)

    err = d.DialAndSend(m)
    return err // Return the error to the consumer loop
}