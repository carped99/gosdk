package outbox

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMessage_Validate(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
		wantErr bool
	}{
		{
			name: "valid message",
			message: &Message{
				EventID:       uuid.New(),
				EventTopic:    "test.topic",
				EventDomain:   "test.domain",
				EventType:     "test.type",
				ObjectType:    "test.object",
				Producer:      "test.producer",
				CorrelationID: "test.correlation",
				Payload:       json.RawMessage(`{"key": "value"}`),
				Metadata:      json.RawMessage(`{"meta": "data"}`),
				CreatedAt:     time.Now().UTC(),
			},
			wantErr: false,
		},
		{
			name: "empty event topic",
			message: &Message{
				EventID:     uuid.New(),
				EventDomain: "test.domain",
				EventType:   "test.type",
				ObjectType:  "test.object",
				CreatedAt:   time.Now().UTC(),
			},
			wantErr: true,
		},
		{
			name: "empty event domain",
			message: &Message{
				EventID:    uuid.New(),
				EventTopic: "test.topic",
				EventType:  "test.type",
				ObjectType: "test.object",
				CreatedAt:  time.Now().UTC(),
			},
			wantErr: true,
		},
		{
			name: "empty event type",
			message: &Message{
				EventID:     uuid.New(),
				EventTopic:  "test.topic",
				EventDomain: "test.domain",
				ObjectType:  "test.object",
				CreatedAt:   time.Now().UTC(),
			},
			wantErr: true,
		},
		{
			name: "empty object type",
			message: &Message{
				EventID:     uuid.New(),
				EventTopic:  "test.topic",
				EventDomain: "test.domain",
				EventType:   "test.type",
				CreatedAt:   time.Now().UTC(),
			},
			wantErr: true,
		},
		{
			name: "empty timestamp",
			message: &Message{
				EventID:     uuid.New(),
				EventTopic:  "test.topic",
				EventDomain: "test.domain",
				EventType:   "test.type",
				ObjectType:  "test.object",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.message.validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
