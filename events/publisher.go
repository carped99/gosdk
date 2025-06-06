package events

import (
	"context"
)

type Publisher interface {
	Publish(ctx context.Context, messages ...*Message) error

	Close() error
}
