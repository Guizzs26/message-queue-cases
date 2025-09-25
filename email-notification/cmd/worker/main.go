package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Guizzs26/message-queue/internal/core/notification"
	"github.com/Guizzs26/message-queue/internal/platform/broker"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RabbitMQURL string `envconfig:"RABBITMQ_URL" required:"true"`
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

	gracefulShutdown(rmqBroker)
}

func run(b broker.Broker) {
	log.Println("Notification worker starting...")

	msgs, err := b.Consume(context.Background(), "email_notifications")
	if err != nil {
		log.Fatalf("failed to start consuming messages: %v", err)
	}

	log.Println("Worker is alive and waiting for messages...")

	var payload notification.EmailPayload
	for msg := range msgs {
		err = json.Unmarshal(msg.Body, &payload)
		if err != nil {
			log.Printf("[ERROR]: failed to unmarshal message body: %v", err)

			if err := msg.Nack(false, false); err != nil {
				log.Printf("[ERROR]: failed to nack message: %v", err)
			}
			continue
		}

		log.Printf("Simulating sending email to %s with subject '%s'", payload.To, payload.Subject)
		time.Sleep(1 * time.Second)

		if err := msg.Ack(false); err != nil {
			log.Printf("[ERROR]: failed to ack message: %v", err)
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

func gracefulShutdown(b broker.Broker) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go run(b)

	<-sigChan
	log.Println("Shutdown signal received")

	// Fecha a conexão com o broker
	b.Close()
	log.Println("Broker connection closed. Shutdown completed.")
}
