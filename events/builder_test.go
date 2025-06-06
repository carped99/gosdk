package events

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessageBuilder(t *testing.T) {
	t.Run("build valid message", func(t *testing.T) {
		// Given
		now := time.Now().UTC()
		id := uuid.New()
		payload := json.RawMessage(`{"key": "value"}`)
		metadata := json.RawMessage(`{"meta": "data"}`)

		// When
		message, err := NewMessageBuilder().
			SetUUID(id).
			SetCreatedAt(now).
			SetEventTopic("test.topic").
			SetEventDomain("test.domain").
			SetEventType("test.type").
			SetObjectType("test.object").
			SetProducer("test.producer").
			SetCorrelationID("test.correlation").
			SetPayload(payload).
			SetMetadata(metadata).
			Build()

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.Equal(t, id, message.EventID)
		assert.Equal(t, now, message.CreatedAt)
		assert.Equal(t, "test.topic", message.EventTopic)
		assert.Equal(t, "test.domain", message.EventDomain)
		assert.Equal(t, "test.type", message.EventType)
		assert.Equal(t, "test.object", message.ObjectType)
		assert.Equal(t, "test.producer", message.Producer)
		assert.Equal(t, "test.correlation", message.CorrelationID)
		assert.Equal(t, payload, message.Payload)
		assert.Equal(t, metadata, message.Metadata)
	})

	t.Run("build invalid message", func(t *testing.T) {
		// Given
		builder := NewMessageBuilder()

		// When
		message, err := builder.Build()

		// Then
		assert.Error(t, err)
		assert.Nil(t, message)
	})

	t.Run("build with default values", func(t *testing.T) {
		// When
		message, err := NewMessageBuilder().
			SetEventTopic("test.topic").
			SetEventDomain("test.domain").
			SetEventType("test.type").
			SetObjectType("test.object").
			Build()

		// Then
		assert.NoError(t, err)
		assert.NotNil(t, message)
		assert.NotEqual(t, uuid.Nil, message.EventID)
		assert.False(t, message.CreatedAt.IsZero())
		assert.Equal(t, "test.topic", message.EventTopic)
		assert.Equal(t, "test.domain", message.EventDomain)
		assert.Equal(t, "test.type", message.EventType)
		assert.Equal(t, "test.object", message.ObjectType)
		assert.Empty(t, message.Producer)
		assert.Empty(t, message.CorrelationID)
		assert.Empty(t, message.Payload)
		assert.Empty(t, message.Metadata)
	})
}

func TestMessageBuilder_Chaining(t *testing.T) {
	// Given
	builder := NewMessageBuilder()

	// When
	builder.
		SetEventTopic("test.topic").
		SetEventDomain("test.domain").
		SetEventType("test.type").
		SetObjectType("test.object")

	// Then
	assert.Equal(t, "test.topic", builder.message.EventTopic)
	assert.Equal(t, "test.domain", builder.message.EventDomain)
	assert.Equal(t, "test.type", builder.message.EventType)
	assert.Equal(t, "test.object", builder.message.ObjectType)
}
