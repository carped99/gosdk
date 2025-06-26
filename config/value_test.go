package config

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAtomicValue_Bool(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected bool
		wantErr  bool
	}{
		{
			name:     "bool true",
			input:    true,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string true",
			input:    "true",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string false",
			input:    "false",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "string 1",
			input:    "1",
			expected: true,
			wantErr:  false,
		},
		{
			name:     "string 0",
			input:    "0",
			expected: false,
			wantErr:  false,
		},
		{
			name:     "int 1",
			input:    1,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "int 0",
			input:    0,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "float 1.0",
			input:    1.0,
			expected: true,
			wantErr:  false,
		},
		{
			name:     "float 0.0",
			input:    0.0,
			expected: false,
			wantErr:  false,
		},
		{
			name:     "invalid string",
			input:    "invalid",
			expected: false,
			wantErr:  true,
		},
		{
			name:     "unsupported type",
			input:    []string{"test"},
			expected: false,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Bool()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicValue_Int(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected int64
		wantErr  bool
	}{
		{
			name:     "int",
			input:    42,
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "int8",
			input:    int8(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "int16",
			input:    int16(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "int32",
			input:    int32(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    int64(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint8",
			input:    uint8(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint16",
			input:    uint16(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint32",
			input:    uint32(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    uint64(42),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    float32(42.5),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    float64(42.5),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "string valid",
			input:    "42",
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "string invalid",
			input:    "invalid",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "unsupported type",
			input:    []int{42},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Int()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicValue_Float(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected float64
		wantErr  bool
	}{
		{
			name:     "int",
			input:    42,
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "int8",
			input:    int8(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "int16",
			input:    int16(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "int32",
			input:    int32(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    int64(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "uint8",
			input:    uint8(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "uint16",
			input:    uint16(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "uint32",
			input:    uint32(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    uint64(42),
			expected: 42.0,
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    float32(42.5),
			expected: 42.5,
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    float64(42.5),
			expected: 42.5,
			wantErr:  false,
		},
		{
			name:     "string valid",
			input:    "42.5",
			expected: 42.5,
			wantErr:  false,
		},
		{
			name:     "string invalid",
			input:    "invalid",
			expected: 0.0,
			wantErr:  true,
		},
		{
			name:     "unsupported type",
			input:    []float64{42.5},
			expected: 0.0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Float()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
		wantErr  bool
	}{
		{
			name:     "string",
			input:    "hello",
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "bool true",
			input:    true,
			expected: "true",
			wantErr:  false,
		},
		{
			name:     "bool false",
			input:    false,
			expected: "false",
			wantErr:  false,
		},
		{
			name:     "int",
			input:    42,
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int8",
			input:    int8(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int16",
			input:    int16(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int32",
			input:    int32(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "int64",
			input:    int64(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "uint",
			input:    uint(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "uint8",
			input:    uint8(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "uint16",
			input:    uint16(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "uint32",
			input:    uint32(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "uint64",
			input:    uint64(42),
			expected: "42",
			wantErr:  false,
		},
		{
			name:     "float32",
			input:    float32(42.5),
			expected: "42.5",
			wantErr:  false,
		},
		{
			name:     "float64",
			input:    float64(42.5),
			expected: "42.5",
			wantErr:  false,
		},
		{
			name:     "[]byte",
			input:    []byte("hello"),
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "unsupported type",
			input:    []string{"hello"},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.String()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicValue_Duration(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected time.Duration
		wantErr  bool
	}{
		{
			name:     "time.Duration type",
			input:    time.Hour + 30*time.Minute,
			expected: time.Hour + 30*time.Minute,
			wantErr:  false,
		},
		{
			name:     "string duration - hours and minutes",
			input:    "1h30m",
			expected: time.Hour + 30*time.Minute,
			wantErr:  false,
		},
		{
			name:     "string duration - seconds",
			input:    "30s",
			expected: 30 * time.Second,
			wantErr:  false,
		},
		{
			name:     "string duration - milliseconds",
			input:    "500ms",
			expected: 500 * time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "string duration - complex",
			input:    "2h15m30s500ms",
			expected: 2*time.Hour + 15*time.Minute + 30*time.Second + 500*time.Millisecond,
			wantErr:  false,
		},
		{
			name:     "string duration - decimal hours",
			input:    "1.5h",
			expected: 90 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "integer nanoseconds",
			input:    int64(3600000000000), // 1 hour in nanoseconds
			expected: time.Hour,
			wantErr:  false,
		},
		{
			name:     "float nanoseconds",
			input:    float64(1800000000000), // 30 minutes in nanoseconds
			expected: 30 * time.Minute,
			wantErr:  false,
		},
		{
			name:     "invalid string duration",
			input:    "invalid",
			expected: 0,
			wantErr:  true,
		},
		{
			name:     "unsupported type",
			input:    []string{"test"},
			expected: 0,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Duration()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAtomicValue_Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []Value
		wantErr  bool
	}{
		{
			name:     "[]any with mixed types",
			input:    []any{"hello", 42, true},
			expected: nil, // Will be checked individually
			wantErr:  false,
		},
		{
			name:     "[]any with strings",
			input:    []any{"a", "b", "c"},
			expected: nil, // Will be checked individually
			wantErr:  false,
		},
		{
			name:     "[]any with numbers",
			input:    []any{1, 2, 3},
			expected: nil, // Will be checked individually
			wantErr:  false,
		},
		{
			name:     "not a slice",
			input:    "not a slice",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "[]string (unsupported)",
			input:    []string{"a", "b", "c"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "[]int (unsupported)",
			input:    []int{1, 2, 3},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Slice()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)

			// Check individual elements for supported types
			if slice, ok := tt.input.([]any); ok {
				assert.Equal(t, len(slice), len(result))
				for i, expectedVal := range slice {
					actualVal := result[i].Load()
					assert.Equal(t, expectedVal, actualVal)
				}
			}
		})
	}
}

func TestAtomicValue_Map(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected map[string]Value
		wantErr  bool
	}{
		{
			name:     "map[string]any with mixed types",
			input:    map[string]any{"key1": "value1", "key2": 42, "key3": true},
			expected: nil, // Will be checked individually
			wantErr:  false,
		},
		{
			name:     "map[string]any with strings",
			input:    map[string]any{"a": "1", "b": "2", "c": "3"},
			expected: nil, // Will be checked individually
			wantErr:  false,
		},
		{
			name:     "not a map",
			input:    "not a map",
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "map[string]string (unsupported)",
			input:    map[string]string{"a": "1", "b": "2"},
			expected: nil,
			wantErr:  true,
		},
		{
			name:     "map[int]any (unsupported)",
			input:    map[int]any{1: "a", 2: "b"},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			result, err := v.Map()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, result)

			// Check individual elements for supported types
			if m, ok := tt.input.(map[string]any); ok {
				assert.Equal(t, len(m), len(result))
				for expectedKey, expectedVal := range m {
					actualVal, exists := result[expectedKey]
					assert.True(t, exists)
					val := actualVal.Load()
					assert.Equal(t, expectedVal, val)
				}
			}
		})
	}
}

func TestAtomicValue_Scan(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		target   any
		expected any
		wantErr  bool
	}{
		{
			name:     "scan string to string",
			input:    "hello",
			target:   new(string),
			expected: "hello",
			wantErr:  false,
		},
		{
			name:     "scan int to int",
			input:    42,
			target:   new(int),
			expected: 42,
			wantErr:  false,
		},
		{
			name:     "scan float to float",
			input:    42.5,
			target:   new(float64),
			expected: 42.5,
			wantErr:  false,
		},
		{
			name:     "scan bool to bool",
			input:    true,
			target:   new(bool),
			expected: true,
			wantErr:  false,
		},
		{
			name:     "scan slice to slice",
			input:    []any{"a", "b", "c"},
			target:   new([]any),
			expected: []any{"a", "b", "c"},
			wantErr:  false,
		},
		{
			name:     "scan map to map",
			input:    map[string]any{"key": "value"},
			target:   new(map[string]any),
			expected: map[string]any{"key": "value"},
			wantErr:  false,
		},
		{
			name:     "scan to nil target",
			input:    "hello",
			target:   nil,
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &atomicValue{}
			v.Store(tt.input)

			err := v.Scan(tt.target)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)

			// Dereference the pointer to get the actual value
			val := reflect.ValueOf(tt.target).Elem().Interface()
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestAtomicValue_LoadAndStore(t *testing.T) {
	t.Run("load and store string", func(t *testing.T) {
		v := &atomicValue{}
		expected := "test value"

		v.Store(expected)
		result := v.Load()

		assert.Equal(t, expected, result)
	})

	t.Run("load and store int", func(t *testing.T) {
		v := &atomicValue{}
		expected := 42

		v.Store(expected)
		result := v.Load()

		assert.Equal(t, expected, result)
	})

	t.Run("load and store complex type", func(t *testing.T) {
		v := &atomicValue{}
		expected := map[string]any{"key": "value"}

		v.Store(expected)
		result := v.Load()

		assert.Equal(t, expected, result)
	})

	t.Run("load and store slice", func(t *testing.T) {
		v := &atomicValue{}
		expected := []any{"a", "b", "c"}

		v.Store(expected)
		result := v.Load()

		assert.Equal(t, expected, result)
	})
}

func TestErrValue(t *testing.T) {
	err := assert.AnError
	v := errValue{err: err}

	t.Run("errValue Bool", func(t *testing.T) {
		result, resultErr := v.Bool()
		assert.Equal(t, false, result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Int", func(t *testing.T) {
		result, resultErr := v.Int()
		assert.Equal(t, int64(0), result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Float", func(t *testing.T) {
		result, resultErr := v.Float()
		assert.Equal(t, 0.0, result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue String", func(t *testing.T) {
		result, resultErr := v.String()
		assert.Equal(t, "", result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Duration", func(t *testing.T) {
		result, resultErr := v.Duration()
		assert.Equal(t, time.Duration(0), result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Slice", func(t *testing.T) {
		result, resultErr := v.Slice()
		assert.Nil(t, result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Map", func(t *testing.T) {
		result, resultErr := v.Map()
		assert.Nil(t, result)
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Scan", func(t *testing.T) {
		resultErr := v.Scan("target")
		assert.Equal(t, err, resultErr)
	})

	t.Run("errValue Load", func(t *testing.T) {
		result := v.Load()
		assert.Nil(t, result)
	})

	t.Run("errValue Store", func(t *testing.T) {
		// Store should not panic
		assert.NotPanics(t, func() {
			v.Store("value")
		})
	})
}

func TestAtomicValue_Duration_EdgeCases(t *testing.T) {
	t.Run("zero duration", func(t *testing.T) {
		v := &atomicValue{}
		v.Store("0s")

		result, err := v.Duration()
		require.NoError(t, err)
		assert.Equal(t, time.Duration(0), result)
	})

	t.Run("very large duration", func(t *testing.T) {
		v := &atomicValue{}
		v.Store("8760h") // 1 year

		result, err := v.Duration()
		require.NoError(t, err)
		assert.Equal(t, 8760*time.Hour, result)
	})

	t.Run("negative duration", func(t *testing.T) {
		v := &atomicValue{}
		v.Store("-1h")

		result, err := v.Duration()
		require.NoError(t, err)
		assert.Equal(t, -time.Hour, result)
	})
}

func TestAtomicValue_TypeAssertError(t *testing.T) {
	t.Run("type assert error message", func(t *testing.T) {
		v := &atomicValue{}
		v.Store("test")

		// Try to get int from string
		_, err := v.Int()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "strconv.ParseInt")
	})
}
