package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockExecutor struct {
	mock.Mock
	execFunc func(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

func (m *MockExecutor) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if m.execFunc != nil {
		return m.execFunc(ctx, query, args...)
	}
	callArgs := m.Called(ctx, query, args)
	if callArgs.Get(0) == nil {
		return nil, callArgs.Error(1)
	}
	return callArgs.Get(0).(sql.Result), callArgs.Error(1)
}

// MockResult is a mock implementation of sql.Result for testing
type MockResult struct{}

func (m *MockResult) LastInsertId() (int64, error) {
	return 1, nil
}

func (m *MockResult) RowsAffected() (int64, error) {
	return 1, nil
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
	publisher, err := NewPublisher(executor)
	require.NoError(t, err)
	defer publisher.Close()

	msg := &Message{
		EventID:       "test-event-1",
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

	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(&MockResult{}, nil)

	// When
	err = publisher.Publish(context.Background(), msg)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_EmptyMessages(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	publisher, err := NewPublisher(executor)
	require.NoError(t, err)
	defer publisher.Close()

	// When
	err = publisher.Publish(context.Background())

	// Then
	assert.NoError(t, err)
	executor.AssertNotCalled(t, "ExecContext")
}

func TestPublisher_Publish_WithRetry(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	publisher, err := NewPublisher(executor, WithMaxRetries(2))
	require.NoError(t, err)
	defer publisher.Close()

	msg := &Message{
		EventID:       "test-event-1",
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
		Return(&MockResult{}, nil).Once()

	// When
	err = publisher.Publish(context.Background(), msg)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_MaxRetriesExceeded(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	publisher, err := NewPublisher(executor, WithMaxRetries(2))
	require.NoError(t, err)
	defer publisher.Close()

	msg := &Message{
		EventID:       "test-event-1",
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
	err = publisher.Publish(context.Background(), msg)

	// Then
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish message")
	executor.AssertExpectations(t)
}

func TestPublisher_Publish_Batch(t *testing.T) {
	// Given
	executor := new(MockExecutor)
	publisher, err := NewPublisher(executor, WithBatchSize(2))
	require.NoError(t, err)
	defer publisher.Close()

	messages := []*Message{
		{
			EventID:       "test-event-1",
			EventTopic:    "test-topic",
			EventDomain:   "test-domain",
			EventType:     "test-event",
			ObjectType:    "test-object",
			Producer:      "test-producer",
			CorrelationID: "test-correlation",
			Payload:       json.RawMessage(`{"key": "value1"}`),
			Metadata:      json.RawMessage(`{"meta": "data1"}`),
			CreatedAt:     time.Now(),
		},
		{
			EventID:       "test-event-2",
			EventTopic:    "test-topic",
			EventDomain:   "test-domain",
			EventType:     "test-event",
			ObjectType:    "test-object",
			Producer:      "test-producer",
			CorrelationID: "test-correlation",
			Payload:       json.RawMessage(`{"key": "value2"}`),
			Metadata:      json.RawMessage(`{"meta": "data2"}`),
			CreatedAt:     time.Now(),
		},
		{
			EventID:       "test-event-3",
			EventTopic:    "test-topic",
			EventDomain:   "test-domain",
			EventType:     "test-event",
			ObjectType:    "test-object",
			Producer:      "test-producer",
			CorrelationID: "test-correlation",
			Payload:       json.RawMessage(`{"key": "value3"}`),
			Metadata:      json.RawMessage(`{"meta": "data3"}`),
			CreatedAt:     time.Now(),
		},
	}

	executor.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(&MockResult{}, nil).Times(3)

	// When
	err = publisher.Publish(context.Background(), messages...)

	// Then
	assert.NoError(t, err)
	executor.AssertExpectations(t)
}

func TestPublisherWithParameterizedQuery(t *testing.T) {
	mockDB := &MockExecutor{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			// Verify the query uses parameterized placeholders
			if !strings.Contains(query, "?") {
				t.Errorf("Expected parameterized query with '?' placeholders, got: %s", query)
			}
			return &MockResult{}, nil
		},
	}

	pub, err := NewPublisher(mockDB)
	require.NoError(t, err)
	defer pub.Close()

	msg := &Message{
		EventID:       "test-event-1",
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-type",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"test": "data"}`),
		Metadata:      json.RawMessage(`{"key": "value"}`),
		CreatedAt:     time.Now(),
	}

	ctx := context.Background()
	err = pub.Publish(ctx, msg)
	require.NoError(t, err)
}

func TestPublisherQueryReuse(t *testing.T) {
	execCount := 0
	mockDB := &MockExecutor{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			execCount++
			return &MockResult{}, nil
		},
	}

	pub, err := NewPublisher(mockDB)
	require.NoError(t, err)
	defer pub.Close()

	msg := &Message{
		EventID:       "test-event-1",
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-type",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"test": "data"}`),
		Metadata:      json.RawMessage(`{"key": "value"}`),
		CreatedAt:     time.Now(),
	}

	ctx := context.Background()

	// First publish should execute the query
	err = pub.Publish(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, 1, execCount)

	// Second publish should execute the query again
	err = pub.Publish(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, 2, execCount)
}

func TestPublisherExecutionError(t *testing.T) {
	mockDB := &MockExecutor{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return nil, fmt.Errorf("execution failed")
		},
	}

	pub, err := NewPublisher(mockDB)
	require.NoError(t, err)
	defer pub.Close()

	msg := &Message{
		EventID:       "test-event-1",
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-type",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"test": "data"}`),
		Metadata:      json.RawMessage(`{"key": "value"}`),
		CreatedAt:     time.Now(),
	}

	ctx := context.Background()
	err = pub.Publish(ctx, msg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to publish message")
}

func TestPublisherClose(t *testing.T) {
	mockDB := &MockExecutor{
		execFunc: func(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
			return &MockResult{}, nil
		},
	}

	pub, err := NewPublisher(mockDB)
	require.NoError(t, err)

	// Publish to create the query
	msg := &Message{
		EventID:       "test-event-1",
		EventTopic:    "test-topic",
		EventDomain:   "test-domain",
		EventType:     "test-type",
		ObjectType:    "test-object",
		Producer:      "test-producer",
		CorrelationID: "test-correlation",
		Payload:       json.RawMessage(`{"test": "data"}`),
		Metadata:      json.RawMessage(`{"key": "value"}`),
		CreatedAt:     time.Now(),
	}

	ctx := context.Background()
	err = pub.Publish(ctx, msg)
	require.NoError(t, err)

	// Close should not error
	err = pub.Close()
	require.NoError(t, err)
}
