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

func NewResendNotifier(k, fe string) Notifier {
	c := resend.NewClient(k)
	return &ResendNotifier{
		client:    c,
		fromEmail: fe,
	}
}

func (rn *ResendNotifier) Send(ctx context.Context, payload []byte) error {
	var ep notification.EmailPayload
	err := json.Unmarshal(payload, &ep)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload for resend: %v", err)
	}

	params := &resend.SendEmailRequest{
		To:      []string{ep.To},
		From:    rn.fromEmail,
		Subject: ep.Subject,
		Html:    ep.Body,
	}

	_, err = rn.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to send email via resend: %v", err)
	}

	return nil
}
