package broker

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

var _ Broker = (*RabbitMQBroker)(nil)

type Message struct {
	ContentType string
	Body        []byte
}

type Broker interface {
	Publish(ctx context.Context, queueName string, msg Message) error
	Close()
}

type RabbitMQBroker struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func NewRabbitMQBroker(url string) (*RabbitMQBroker, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to create rabbitmq channel: %v", err)
	}

	return &RabbitMQBroker{
		conn:    conn,
		channel: channel,
	}, nil
}

func (r *RabbitMQBroker) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.channel != nil {
		r.channel.Close()
	}
}

func (r *RabbitMQBroker) Publish(ctx context.Context, queueName string, msg Message) error {
	_, err := r.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare the rabbitmq queue: %v", err)
	}

	var amqpmsg *amqp.Publishing
	amqpmsg.ContentType = msg.ContentType
	amqpmsg.Body = msg.Body
	amqpmsg.DeliveryMode = amqp.Persistent // disk storage

	err = r.channel.PublishWithContext(
		ctx,
		"",
		queueName,
		false,
		false,
		*amqpmsg,
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}
