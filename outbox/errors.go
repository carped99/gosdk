package outbox

import (
	"errors"
	"fmt"
)

var (
	// ErrInvalidEventTopic is returned when the event topic is empty.
	// The topic is required to determine where the event should be published.
	ErrInvalidEventTopic = errors.New("event topic cannot be empty")

	// ErrInvalidEventDomain is returned when the event domain is not configured.
	// The domain provides context about which part of the system the event belongs to.
	ErrInvalidEventDomain = errors.New("events domain not configured")

	// ErrInvalidEventType is returned when the event type is empty.
	// The event type describes what kind of event occurred (e.g., created, updated, deleted).
	ErrInvalidEventType = errors.New("event type cannot be empty")

	// ErrInvalidObjectType is returned when the object type is empty.
	// The object type identifies what kind of object the event is about.
	ErrInvalidObjectType = errors.New("object type cannot be empty")

	// ErrInvalidTimestamp is returned when the creation timestamp is not set.
	// The timestamp is required to track when the event occurred.
	ErrInvalidTimestamp = errors.New("timestamp cannot be empty")

	ErrInvalidConfig    = errors.New("invalid publisher configuration")
	ErrPublishFailed    = errors.New("failed to publish message")
	ErrMarshalFailed    = errors.New("failed to marshal message")
	ErrExecutorNotFound = errors.New("executor not found")
)

type PublishError struct {
	MessageID string
	Err       error
	Retries   int
}

func (e *PublishError) Error() string {
	return fmt.Sprintf("failed to publish message %s after %d retries: %v", e.MessageID, e.Retries, e.Err)
}

func (e *PublishError) Unwrap() error {
	return e.Err
}

type MarshalError struct {
	Field string
	Err   error
}

func (e *MarshalError) Error() string {
	return fmt.Sprintf("failed to marshal %s: %v", e.Field, e.Err)
}

func (e *MarshalError) Unwrap() error {
	return e.Err
}

type ConfigError struct {
	Field string
	Err   error
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("invalid configuration for %s: %v", e.Field, e.Err)
}

func (e *ConfigError) Unwrap() error {
	return e.Err
}
