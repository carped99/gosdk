package events

import "errors"

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
)
