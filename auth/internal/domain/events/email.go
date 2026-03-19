package events

import "context"

type EmailWorker interface {
	Start(ctx context.Context)
	Stop()
}
