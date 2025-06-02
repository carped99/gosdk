package outbox

import (
	"fmt"
	"github.com/google/uuid"
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

type MessageBuilder struct {
	message *Message
}

func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		message: &Message{
			UUID:      uuid.New(),
			CreatedAt: time.Now().UTC(),
		},
	}
}

func (b *MessageBuilder) SetEventTopic(eventTopic string) *MessageBuilder {
	b.message.EventTopic = eventTopic
	return b
}

func (b *MessageBuilder) SetEventDomain(eventDomain string) *MessageBuilder {
	b.message.EventDomain = eventDomain
	return b
}

func (b *MessageBuilder) SetObjectType(objectType string) *MessageBuilder {
	b.message.ObjectType = objectType
	return b
}

func (b *MessageBuilder) SetEventType(eventType string) *MessageBuilder {
	b.message.EventType = eventType
	return b
}

func (b *MessageBuilder) SetProducer(producer string) *MessageBuilder {
	b.message.Producer = producer
	return b
}

func (b *MessageBuilder) SetCorrelationID(correlationID string) *MessageBuilder {
	b.message.CorrelationID = correlationID
	return b
}

func (b *MessageBuilder) SetPayload(payload any) *MessageBuilder {
	b.message.Payload = payload
	return b
}

func (b *MessageBuilder) SetMetadata(metadata any) *MessageBuilder {
	b.message.Metadata = metadata
	return b
}

func (b *MessageBuilder) Build() (*Message, error) {
	if err := b.message.Validate(); err != nil {
		return nil, err
	}
	return b.message, nil
}

func (m *Message) Validate() error {
	if m.EventTopic == "" {
		return fmt.Errorf("event topic cannot be empty")
	}

	if m.EventDomain == "" {
		return fmt.Errorf("event domain cannot be empty")
	}

	if m.ObjectType == "" {
		return fmt.Errorf("object type cannot be empty")
	}

	if m.EventType == "" {
		return fmt.Errorf("event type cannot be empty")
	}

	if m.CreatedAt.IsZero() {
		return fmt.Errorf("created at cannot be zero")
	}

	return nil
}
