package outbox

import (
	"errors"
	"fmt"
)

var (
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
