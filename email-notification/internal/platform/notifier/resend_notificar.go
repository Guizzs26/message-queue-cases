package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Guizzs26/message-queue/internal/core/notification"
	"github.com/resend/resend-go/v2"
)

type ResendNotifier struct {
	client    *resend.Client
	fromEmail string
}

func NewResendNotifier(apiKey, fromEmail string) Notifier {
	client := resend.NewClient(apiKey)

	return &ResendNotifier{
		client:    client,
		fromEmail: fromEmail,
	}
}

func (rn *ResendNotifier) Send(ctx context.Context, payload []byte) error {
	var data notification.EmailPayload
	err := json.Unmarshal(payload, &data)
	if err != nil {
		return fmt.Errorf("failed to unmarshal email payload to send send email via resend: %v", err)
	}

	params := &resend.SendEmailRequest{
		To:      []string{data.To},
		From:    rn.fromEmail,
		Subject: data.Subject,
		Html:    data.Body,
	}

	_, err = rn.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %v", err)
	}

	return nil
}
