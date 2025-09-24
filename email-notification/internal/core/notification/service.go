package notification

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Guizzs26/message-queue/internal/platform/broker"
)

type EmailService interface {
	SendEmail(ctx context.Context, to, subject, body string) error
}

type EmailPayload struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type NotificationService struct {
	broker broker.Broker
}

func NewNotificationService(b broker.Broker) EmailService {
	return &NotificationService{broker: b}
}

func (ns *NotificationService) SendEmail(
	ctx context.Context,
	to,
	subject,
	body string,
) error {
	p := EmailPayload{
		To:      to,
		Subject: subject,
		Body:    body,
	}

	b, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("marshal email payload: %v", err)
	}

	msg := broker.Message{
		ContentType: "application/json",
		Body:        b,
	}

	if err := ns.broker.Publish(ctx, "email_notifications", msg); err != nil {
		return fmt.Errorf("send email: %v", err)
	}

	return nil
}
