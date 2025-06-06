package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"go.uber.org/multierr"
)

// Message represents an event message in the system.
// It contains all necessary information about an event including its metadata and payload.
type Message struct {
	EventID       uuid.UUID       `json:"event_id"`                 // Unique identifier for the message
	EventTopic    string          `json:"event_topic"`              // Topic where the event will be published
	EventDomain   string          `json:"event_domain"`             // Domain context of the event
	EventType     string          `json:"event_type"`               // Type of the event (e.g., created, updated, deleted)
	ObjectType    string          `json:"object_type"`              // Type of the object this event is about
	Producer      string          `json:"producer,omitempty"`       // Service or component that produced the event
	CorrelationID string          `json:"correlation_id,omitempty"` // ID for tracking related events
	Payload       json.RawMessage `json:"payload,omitempty"`        // Event data in JSON format
	Metadata      json.RawMessage `json:"metadata,omitempty"`       // Additional event information
	CreatedAt     time.Time       `json:"created_at"`               // When the event was created
}

// GetEventID returns the unique identifier of the message
func (m *Message) GetEventID() uuid.UUID {
	return m.EventID
}

// GetEventTopic returns the topic where the event will be published
func (m *Message) GetEventTopic() string {
	return m.EventTopic
}

// GetEventDomain returns the domain context of the event
func (m *Message) GetEventDomain() string {
	return m.EventDomain
}

// GetEventType returns the type of the event
func (m *Message) GetEventType() string {
	return m.EventType
}

// GetObjectType returns the type of object this event is about
func (m *Message) GetObjectType() string {
	return m.ObjectType
}

// GetProducer returns the service or component that produced the event
func (m *Message) GetProducer() string {
	return m.Producer
}

// GetCorrelationID returns the ID for tracking related events
func (m *Message) GetCorrelationID() string {
	return m.CorrelationID
}

// GetPayload returns the event data
func (m *Message) GetPayload() json.RawMessage {
	return m.Payload
}

// GetMetadata returns additional event information
func (m *Message) GetMetadata() json.RawMessage {
	return m.Metadata
}

// GetCreatedAt returns when the event was created
func (m *Message) GetCreatedAt() time.Time {
	return m.CreatedAt
}

// validate checks if all required fields of the Message are properly set.
func (m *Message) validate() error {
	var err error

	if m.EventTopic == "" {
		err = multierr.Append(err, ErrInvalidEventTopic)
	}

	if m.EventDomain == "" {
		err = multierr.Append(err, ErrInvalidEventDomain)
	}

	if m.EventType == "" {
		err = multierr.Append(err, ErrInvalidEventType)
	}

	if m.ObjectType == "" {
		err = multierr.Append(err, ErrInvalidObjectType)
	}

	if m.CreatedAt.IsZero() {
		err = multierr.Append(err, ErrInvalidTimestamp)
	}

	return err
}
