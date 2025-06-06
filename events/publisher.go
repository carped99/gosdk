package events

import (
	"context"
)

// Publisher defines the interface for publishing event messages.
// Implementations of this interface are responsible for delivering events
// to their intended destinations (e.g., message brokers, databases).
type Publisher interface {
	// Publish sends one or more event messages.
	// Returns an error if the publishing operation fails.
	// The error should provide details about what went wrong during publishing.
	Publish(ctx context.Context, messages ...*Message) error

	Close() error
}
