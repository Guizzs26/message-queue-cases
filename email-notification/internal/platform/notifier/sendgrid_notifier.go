package notifier

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Guizzs26/message-queue/internal/core/notification"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridNotifier struct {
	client    *sendgrid.Client
	fromName  string
	fromEmail string
}

func NewSendGridNotifier(k, fn, fe string) Notifier {
	c := sendgrid.NewSendClient(k)
	return &SendGridNotifier{
		client:    c,
		fromName:  fn,
		fromEmail: fe,
	}
}

func (sgn *SendGridNotifier) Send(ctx context.Context, payload []byte) error {
	var ep notification.EmailPayload
	err := json.Unmarshal(payload, &ep)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload for sendgrid: %v", err)
	}

	from := mail.NewEmail(sgn.fromName, sgn.fromEmail)
	to := mail.NewEmail(ep.To, ep.To)
	plainTextContent := ep.Body
	htmlContent := ep.Body

	msg := mail.NewSingleEmail(from, ep.Subject, to, plainTextContent, htmlContent)

	res, err := sgn.client.SendWithContext(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to send email via sendgrid: %v", err)
	}

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	return fmt.Errorf("failed to send email, status code %d, body: %s", res.StatusCode, res.Body)
}
