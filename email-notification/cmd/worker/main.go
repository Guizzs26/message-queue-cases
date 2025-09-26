package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Guizzs26/message-queue/internal/platform/broker"
	"github.com/Guizzs26/message-queue/internal/platform/notifier"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RabbitMQURL  string `envconfig:"RABBITMQ_URL" required:"true"`
	ResendAPIKey string `envconfig:"RESEND_API_KEY" required:"true"`
	EmailFrom    string `envconfig:"EMAIL_FROM_ADDRESS" required:"true"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	rmqBroker, err := broker.NewRabbitMQBroker(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("failed to initialize rabbitmq broker: %v", err)
	}

	resendNotifier := notifier.NewResendNotifier(cfg.ResendAPIKey, cfg.EmailFrom)

	gracefulShutdown(rmqBroker, resendNotifier)
}

func run(b broker.Broker, n notifier.Notifier) {
	log.Println("Notification worker starting...")

	queueName := "email_notifications"
	msgs, err := b.Consume(context.Background(), queueName)
	if err != nil {
		log.Fatalf("failed to start consuming from queue '%s': %v", queueName, err)
	}

	log.Println("Worker is alive and waiting for messages...")

	for msg := range msgs {
		log.Printf("Received a message from queue '%s'", queueName)

		if err := n.Send(context.Background(), msg.Body); err != nil {
			log.Printf("ERROR: failed to send notification: %v. Rejecting message.", err)
			if nackErr := msg.Nack(false, false); nackErr != nil {
				log.Printf("FATAL: failed to nack message: %v", nackErr)
			} else {
				log.Printf("Successfully processed message. Sending Ack.")
				if ackErr := msg.Ack(false); ackErr != nil {
					log.Printf("ERROR: failed to ack message: %v", ackErr)
				}
			}
		}
	}

	log.Println("Message channel closed. Worker shutting down...")
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found or error loading it")
	}
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to process environment config: %v", err)
	}
	return &cfg, nil
}

func gracefulShutdown(b broker.Broker, n notifier.Notifier) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go run(b, n)

	<-sigChan
	log.Println("Shutdown signal received")

	// Fecha a conexão com o broker
	b.Close()
	log.Println("Broker connection closed. Shutdown completed.")
}
