package events

import (
	"github.com/google/uuid"
	"go.uber.org/multierr"
	"time"
)

type Message struct {
	UUID          uuid.UUID `json:"uuid"`
	EventTopic    string    `json:"event_topic"`
	EventDomain   string    `json:"event_domain"`
	EventType     string    `json:"event_type"`
	ObjectType    string    `json:"object_type"`
	Producer      string    `json:"producer,omitempty"`
	CorrelationID string    `json:"correlation_id,omitempty"`
	Payload       any       `json:"payload,omitempty"`
	Metadata      any       `json:"metadata,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

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
