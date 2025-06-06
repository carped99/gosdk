package events

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

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

func (b *MessageBuilder) SetUUID(id uuid.UUID) *MessageBuilder {
	b.message.UUID = id
	return b
}

func (b *MessageBuilder) SetCreatedAt(value time.Time) *MessageBuilder {
	b.message.CreatedAt = value
	return b
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

func (b *MessageBuilder) SetPayload(payload json.RawMessage) *MessageBuilder {
	b.message.Payload = payload
	return b
}

func (b *MessageBuilder) SetMetadata(metadata json.RawMessage) *MessageBuilder {
	b.message.Metadata = metadata
	return b
}

func (b *MessageBuilder) Build() (*Message, error) {
	if err := b.message.validate(); err != nil {
		return nil, err
	}
	return b.message, nil
}
