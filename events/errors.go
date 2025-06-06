package events

import "errors"

var (
	ErrInvalidEventTopic  = errors.New("event topic cannot be empty")
	ErrInvalidEventDomain = errors.New("events domain not configured")
	ErrInvalidEventType   = errors.New("event type cannot be empty")
	ErrInvalidObjectType  = errors.New("object type cannot be empty")
	ErrInvalidTimestamp   = errors.New("timestamp cannot be empty")
)
