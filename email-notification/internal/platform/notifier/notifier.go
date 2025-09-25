package notifier

import "context"

type Notifier interface {
	Send(ctx context.Context, payload []byte) error
}
