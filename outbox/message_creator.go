package outbox

import (
	"encoding/json"
	"reflect"

	"github.com/huandu/xstrings"
)

// CreateEventOption defines a function that creates event
type CreateEventOption[T any] func(value T, c *createEventConfig[T])

type createEventConfig[T any] struct {
	eventTopic    string
	eventDomain   string
	eventType     string
	objectType    string
	producer      string
	correlationID string
	payload       json.RawMessage
	metadata      json.RawMessage
}

// WithEventTopic sets the event topic
func WithEventTopic[T any](topic string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.eventTopic = topic
	}
}

// WithEventTopicFunc sets the event topic using a function
func WithEventTopicFunc[T any](fn func(T) string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.eventTopic = fn(value)
	}
}

// WithEventDomainFunc sets the event domain using a function
func WithEventDomainFunc[T any](fn func(T) string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.eventDomain = fn(value)
	}
}

// WithObjectType sets the object type
func WithObjectType[T any](objectType string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.objectType = objectType
	}
}

// WithObjectTypeFunc sets the object type using a function
func WithObjectTypeFunc[T any](fn func(T) string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.objectType = fn(value)
	}
}

// WithPayload sets the payload
func WithPayload[T any](payload json.RawMessage) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.payload = payload
	}
}

// WithPayloadFunc sets the payload using a function
func WithPayloadFunc[T any](fn func(T) json.RawMessage) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.payload = fn(value)
	}
}

// WithProducer sets the producer
func WithProducer[T any](producer string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.producer = producer
	}
}

// WithProducerFunc sets the producer using a function
func WithProducerFunc[T any](fn func(T) string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.producer = fn(value)
	}
}

// WithCorrelationID sets the correlation ID
func WithCorrelationID[T any](correlationID string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.correlationID = correlationID
	}
}

// WithCorrelationIDFunc sets the correlation ID using a function
func WithCorrelationIDFunc[T any](fn func(T) string) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.correlationID = fn(value)
	}
}

// WithMetadata sets the metadata
func WithMetadata[T any](metadata json.RawMessage) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.metadata = metadata
	}
}

// WithMetadataFunc sets the metadata using a function
func WithMetadataFunc[T any](fn func(T) json.RawMessage) CreateEventOption[T] {
	return func(value T, c *createEventConfig[T]) {
		c.metadata = fn(value)
	}
}

// CreateMessage creates a new event message with the given domain, type, and value
func CreateMessage[T any](eventDomain, eventType string, value T, opts ...CreateEventOption[T]) (*Message, error) {
	config := &createEventConfig[T]{
		eventTopic:  eventDomain + ".events",
		eventDomain: eventDomain,
		eventType:   eventType,
	}

	// Apply options
	for _, opt := range opts {
		opt(value, config)
	}

	// Set default objectType if not provided
	if config.objectType == "" {
		config.objectType = getObjectTypeFromValue(value)
	}

	// Set default payload if not provided
	if config.payload == nil {
		payload, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		config.payload = payload
	}

	return NewMessageBuilder().
		SetEventTopic(config.eventTopic).
		SetEventDomain(config.eventDomain).
		SetEventType(config.eventType).
		SetObjectType(config.objectType).
		SetProducer(config.producer).
		SetCorrelationID(config.correlationID).
		SetPayload(config.payload).
		SetMetadata(config.metadata).
		Build()
}

// getObjectTypeFromValue extracts the object type from a value
func getObjectTypeFromValue[T any](v T) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return xstrings.ToSnakeCase(t.Name())
}
