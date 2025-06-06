package events

import (
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
