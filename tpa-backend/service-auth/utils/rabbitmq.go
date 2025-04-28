package utils

import (
	"fmt"
	"log"

	"github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
    if err != nil {
        log.Fatalf("%s: %s", msg, err)
    }
}

func GetRabbitMQChannel() (*amqp091.Channel, error) {
	rabbitURL := GetDotEnv("RABBITMQ_URL") // "amqp://guest:guest@localhost:5672/"
    if rabbitURL == "" {
        log.Fatal("RABBITMQ_URL not set in .env")
    }
	
    conn, err := amqp091.Dial(rabbitURL)
    failOnError(err, "Failed to connect to RabbitMQ")
    log.Println("RabbitMQ connection established")

	ch, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open a channel: %w", err)
    }
    return ch, nil
}