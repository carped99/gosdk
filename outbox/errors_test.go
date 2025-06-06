package outbox

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "ErrInvalidEventTopic",
			err:      ErrInvalidEventTopic,
			expected: "event topic cannot be empty",
		},
		{
			name:     "ErrInvalidEventDomain",
			err:      ErrInvalidEventDomain,
			expected: "events domain not configured",
		},
		{
			name:     "ErrInvalidEventType",
			err:      ErrInvalidEventType,
			expected: "event type cannot be empty",
		},
		{
			name:     "ErrInvalidObjectType",
			err:      ErrInvalidObjectType,
			expected: "object type cannot be empty",
		},
		{
			name:     "ErrInvalidTimestamp",
			err:      ErrInvalidTimestamp,
			expected: "timestamp cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestPublishError(t *testing.T) {
	// Given
	originalErr := errors.New("original error")
	publishErr := &PublishError{
		MessageID: "test-message",
		Err:       originalErr,
		Retries:   3,
	}

	// When
	errMsg := publishErr.Error()
	unwrappedErr := publishErr.Unwrap()

	// Then
	assert.Equal(t, "failed to publish message test-message after 3 retries: original error", errMsg)
	assert.Equal(t, originalErr, unwrappedErr)
}

func TestMarshalError(t *testing.T) {
	// Given
	originalErr := errors.New("json error")
	marshalErr := &MarshalError{
		Field: "payload",
		Err:   originalErr,
	}

	// When
	errMsg := marshalErr.Error()
	unwrappedErr := marshalErr.Unwrap()

	// Then
	assert.Equal(t, "failed to marshal payload: json error", errMsg)
	assert.Equal(t, originalErr, unwrappedErr)
}

func TestConfigError(t *testing.T) {
	// Given
	originalErr := errors.New("invalid value")
	configErr := &ConfigError{
		Field: "MaxRetries",
		Err:   originalErr,
	}

	// When
	errMsg := configErr.Error()
	unwrappedErr := configErr.Unwrap()

	// Then
	assert.Equal(t, "invalid configuration for MaxRetries: invalid value", errMsg)
	assert.Equal(t, originalErr, unwrappedErr)
}

func TestErrorConstants(t *testing.T) {
	assert.Equal(t, "invalid publisher configuration", ErrInvalidConfig.Error())
	assert.Equal(t, "failed to publish message", ErrPublishFailed.Error())
	assert.Equal(t, "failed to marshal message", ErrMarshalFailed.Error())
	assert.Equal(t, "executor not found", ErrExecutorNotFound.Error())
}
