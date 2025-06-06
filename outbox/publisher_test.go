package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockExecutor struct {
	mock.Mock
}

func (m *MockExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	callArgs := m.Called(ctx, query, args)
	return nil, callArgs.Error(1)
}

// testMessage is a test implementation of the Message interface
type testMessage struct {
	UUID          uuid.UUID
	EventTopic    string
	EventDomain   string
	EventType     string
	ObjectType    string
	Producer      string
	CorrelationID string
	Payload       json.RawMessage
	Metadata      json.RawMessage
	CreatedAt     time.Time
}

func (m *testMessage) GetEventID() uuid.UUID {
	return m.UUID
}

func (m *testMessage) GetEventTopic() string {
	return m.EventTopic
}

func (m *testMessage) GetEventDomain() string {
	return m.EventDomain
}

func (m *testMessage) GetEventType() string {
	return m.EventType
}

func (m *testMessage) GetObjectType() string {
	return m.ObjectType
}

func (m *testMessage) GetProducer() string {
	return m.Producer
}

func (m *testMessage) GetCorrelationID() string {
	return m.CorrelationID
}

func (m *testMessage) GetPayload() json.RawMessage {
	return m.Payload
}

func (m *testMessage) GetMetadata() json.RawMessage {
	return m.Metadata
}

func (m *testMessage) GetCreatedAt() time.Time {
	return m.CreatedAt
}

func TestPublisher_Publish_Success(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	config := PublisherConfig{
		BatchSize:  10,
		MaxRetries: 3,
	}
	publisher := NewPublisher(executor, config)

	msg := &testMessage{
		UUID:          uuid.New(),
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-event",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"key": "value"}`),
		Metadata:      json.RawMessage(`{"meta": "data"}`),
		CreatedAt:     time.Now(),
	}

	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

	// When
	err := publisher.Publish(context.Background(), msg)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_EmptyMessages(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	config := PublisherConfig{
		BatchSize:  10,
		MaxRetries: 3,
	}
	publisher := NewPublisher(executor, config)

	// When
	err := publisher.Publish(context.Background())

	// Then
	assert.NoError(t, err)
	executor.AssertNotCalled(t, "ExecContext")
}

func TestPublisher_Publish_WithRetry(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	config := PublisherConfig{
		BatchSize:  10,
		MaxRetries: 2,
	}
	publisher := NewPublisher(executor, config)

	msg := &testMessage{
		UUID:          uuid.New(),
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-event",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"key": "value"}`),
		Metadata:      json.RawMessage(`{"meta": "data"}`),
		CreatedAt:     time.Now(),
	}

	// First attempt fails, second succeeds
	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("db error")).Once()
	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, nil).Once()

	// When
	err := publisher.Publish(context.Background(), msg)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_MaxRetriesExceeded(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	config := PublisherConfig{
		BatchSize:  10,
		MaxRetries: 2,
	}
	publisher := NewPublisher(executor, config)

	msg := &testMessage{
		UUID:          uuid.New(),
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-event",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"key": "value"}`),
		Metadata:      json.RawMessage(`{"meta": "data"}`),
		CreatedAt:     time.Now(),
	}

	// All attempts fail
	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("db error")).Times(3)

	// When
	err := publisher.Publish(context.Background(), msg)

	// Then
	assert.Error(t, err)
	assert.Equal(t, "db error", err.Error())
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_Batch(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	config := PublisherConfig{
		BatchSize:  10,
		MaxRetries: 3,
	}
	publisher := NewPublisher(executor, config)

	messages := []Message{
		&testMessage{
			UUID:          uuid.New(),
			EventTopic:    "test-topic-1",
			EventDomain:   "test-domain",
			EventType:     "test-event",
			ObjectType:    "test-object",
			Producer:      "test-producer",
			CorrelationID: "test-correlation",
			Payload:       json.RawMessage(`{"key": "value1"}`),
			Metadata:      json.RawMessage(`{"meta": "data1"}`),
			CreatedAt:     time.Now(),
		},
		&testMessage{
			UUID:          uuid.New(),
			EventTopic:    "test-topic-2",
			EventDomain:   "test-domain",
			EventType:     "test-event",
			ObjectType:    "test-object",
			Producer:      "test-producer",
			CorrelationID: "test-correlation",
			Payload:       json.RawMessage(`{"key": "value2"}`),
			Metadata:      json.RawMessage(`{"meta": "data2"}`),
			CreatedAt:     time.Now(),
		},
	}

	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil).Times(2)

	// When
	err := publisher.Publish(context.Background(), messages...)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}
