package outbox

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

// MessageBuilder implements the builder pattern for creating Message objects.
// It provides a fluent interface for setting message properties and ensures
// all required fields are properly set before creating the message.
type MessageBuilder struct {
	message *Message
}

// NewMessageBuilder creates a new messageBuilder instance.
// It initializes the message with a new EventID and current timestamp.
func NewMessageBuilder() *MessageBuilder {
	return &MessageBuilder{
		message: &Message{
			CreatedAt: time.Now().UTC(),
		},
	}
}

// SetEventID sets the unique identifier for the message.
func (b *MessageBuilder) SetEventID(id string) *MessageBuilder {
	b.message.EventID = id
	return b
}

// SetEventTopic sets the topic where the event will be published.
func (b *MessageBuilder) SetEventTopic(eventTopic string) *MessageBuilder {
	b.message.EventTopic = eventTopic
	return b
}

// SetEventDomain sets the domain context of the event.
func (b *MessageBuilder) SetEventDomain(eventDomain string) *MessageBuilder {
	b.message.EventDomain = eventDomain
	return b
}

// SetEventType sets the type of event (e.g., created, updated, deleted).
func (b *MessageBuilder) SetEventType(eventType string) *MessageBuilder {
	b.message.EventType = eventType
	return b
}

// SetObjectType sets the type of object this event is about.
func (b *MessageBuilder) SetObjectType(objectType string) *MessageBuilder {
	b.message.ObjectType = objectType
	return b
}

// SetProducer sets the service or component that produced the event.
func (b *MessageBuilder) SetProducer(producer string) *MessageBuilder {
	b.message.Producer = producer
	return b
}

// SetCorrelationID sets the ID for tracking related events.
func (b *MessageBuilder) SetCorrelationID(correlationID string) *MessageBuilder {
	b.message.CorrelationID = correlationID
	return b
}

// SetPayload sets the event data in JSON format.
func (b *MessageBuilder) SetPayload(payload json.RawMessage) *MessageBuilder {
	b.message.Payload = payload
	return b
}

// SetMetadata sets additional information about the event.
func (b *MessageBuilder) SetMetadata(metadata json.RawMessage) *MessageBuilder {
	b.message.Metadata = metadata
	return b
}

// SetCreatedAt sets the timestamp when the message was created.
func (b *MessageBuilder) SetCreatedAt(value time.Time) *MessageBuilder {
	b.message.CreatedAt = value
	return b
}

// Build validates all required fields and creates a new Message.
// Returns an error if any required field is missing or invalid.
func (b *MessageBuilder) Build() (*Message, error) {
	if b.message.EventID == "" {
		b.message.EventID = uuid.New().String()
	}

	if b.message.CreatedAt.IsZero() {
		b.message.CreatedAt = time.Now().UTC()
	}

	if err := b.message.validate(); err != nil {
		return nil, err
	}

	// 메시지 복사본 생성
	message := *b.message
	return &message, nil
}
