package outbox

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Message defines the interface for event messages that can be published through the outbox.
// It provides access to all necessary message properties through getter methods.
type Message interface {
	GetEventID() uuid.UUID
	GetEventTopic() string
	GetEventDomain() string
	GetEventType() string
	GetObjectType() string
	GetProducer() string
	GetCorrelationID() string
	GetPayload() json.RawMessage
	GetMetadata() json.RawMessage
	GetCreatedAt() time.Time
}
