package outbox

import (
	"encoding/json"
	"time"

	"go.uber.org/multierr"
)

// Message represents an event message in the system.
// It contains all necessary information about an event including its metadata and payload.
type Message struct {
	EventID       string          `json:"event_id"`                 // Unique identifier for the message
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

// validate checks if all required fields of the Message are properly set.
func (m *Message) validate() error {
	var err error

	if m.EventID == "" {
		err = multierr.Append(err, ErrInvalidEventID)
	}

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

// UnmarshalMessage unmarshal a JSON byte slice into a Message struct.
func UnmarshalMessage(data []byte) (*Message, error) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}

	if err := msg.validate(); err != nil {
		return nil, err
	}

	return &msg, nil
}
