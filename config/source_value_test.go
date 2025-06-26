package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewValueSourceFromMap(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]any
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple map",
			input: map[string]any{
				"key1": "value1",
				"key2": 123,
				"key3": true,
			},
			want: map[string]any{
				"key1": "value1",
				"key2": 123,
				"key3": true,
			},
			wantErr: false,
		},
		{
			name: "nested map",
			input: map[string]any{
				"database": map[string]any{
					"host": "localhost",
					"port": 5432,
				},
				"redis": map[string]any{
					"host": "127.0.0.1",
					"port": 6379,
				},
			},
			want: map[string]any{
				"database.host": "localhost",
				"database.port": 5432,
				"redis.host":    "127.0.0.1",
				"redis.port":    6379,
			},
			wantErr: false,
		},
		{
			name:    "empty map",
			input:   map[string]any{},
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name:    "nil map",
			input:   nil,
			want:    map[string]any{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, err := NewValueSourceFromMap(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, source)

			result, err := source.Load()
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestNewValueSourceFromStruct(t *testing.T) {
	type TestConfig struct {
		Database struct {
			Host string `koanf:"host"`
			Port int    `koanf:"port"`
		} `koanf:"database"`
		Redis struct {
			Host string `koanf:"host"`
			Port int    `koanf:"port"`
		} `koanf:"redis"`
		Debug bool `koanf:"debug"`
	}

	tests := []struct {
		name    string
		input   any
		tag     string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "valid struct with koanf tags",
			input: TestConfig{
				Database: struct {
					Host string `koanf:"host"`
					Port int    `koanf:"port"`
				}{
					Host: "localhost",
					Port: 5432,
				},
				Redis: struct {
					Host string `koanf:"host"`
					Port int    `koanf:"port"`
				}{
					Host: "127.0.0.1",
					Port: 6379,
				},
				Debug: true,
			},
			tag: "koanf",
			want: map[string]any{
				"database.host": "localhost",
				"database.port": 5432,
				"redis.host":    "127.0.0.1",
				"redis.port":    6379,
				"debug":         true,
			},
			wantErr: false,
		},
		{
			name:    "nil input",
			input:   nil,
			tag:     "koanf",
			want:    map[string]any{},
			wantErr: false,
		},
		{
			name: "struct without tags",
			input: struct {
				Name string
				Age  int
			}{
				Name: "test",
				Age:  25,
			},
			tag: "koanf",
			want: map[string]any{
				"Name": "test",
				"Age":  25,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, err := NewValueSourceFromStruct(tt.input, tt.tag)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, source)

			result, err := source.Load()
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestValueSource_Load(t *testing.T) {
	input := map[string]any{
		"app": map[string]any{
			"name":    "test-app",
			"version": "1.0.0",
		},
		"server": map[string]any{
			"port": 8080,
			"host": "0.0.0.0",
		},
	}

	source, err := NewValueSourceFromMap(input)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	expected := map[string]any{
		"app.name":    "test-app",
		"app.version": "1.0.0",
		"server.port": 8080,
		"server.host": "0.0.0.0",
	}
	assert.Equal(t, expected, result)

	// Load should return the same data multiple times
	result2, err := source.Load()
	require.NoError(t, err)
	assert.Equal(t, expected, result2)
}

func TestValueSource_Watch(t *testing.T) {
	input := map[string]any{
		"key": "value",
	}

	source, err := NewValueSourceFromMap(input)
	require.NoError(t, err)

	watcher, err := source.Watch()
	require.NoError(t, err)
	assert.Nil(t, watcher) // ValueSource의 Watch는 no-op이므로 nil을 반환해야 함
}

func TestValueSource_Integration(t *testing.T) {
	// Map에서 생성한 소스 테스트
	mapInput := map[string]any{
		"config": map[string]any{
			"environment": "test",
			"log_level":   "debug",
		},
	}

	mapSource, err := NewValueSourceFromMap(mapInput)
	require.NoError(t, err)

	mapResult, err := mapSource.Load()
	require.NoError(t, err)
	expectedMapResult := map[string]any{
		"config.environment": "test",
		"config.log_level":   "debug",
	}
	assert.Equal(t, expectedMapResult, mapResult)

	// Struct에서 생성한 소스 테스트
	type Config struct {
		Environment string `koanf:"environment"`
		LogLevel    string `koanf:"log_level"`
	}

	structInput := Config{
		Environment: "production",
		LogLevel:    "info",
	}

	structSource, err := NewValueSourceFromStruct(structInput, "koanf")
	require.NoError(t, err)

	structResult, err := structSource.Load()
	require.NoError(t, err)
	expectedStructResult := map[string]any{
		"environment": "production",
		"log_level":   "info",
	}
	assert.Equal(t, expectedStructResult, structResult)
}

func TestValueSource_ComplexNestedStructure(t *testing.T) {
	complexInput := map[string]any{
		"application": map[string]any{
			"name":    "my-app",
			"version": "2.1.0",
			"features": map[string]any{
				"auth":  true,
				"cache": false,
			},
		},
		"databases": map[string]any{
			"primary": map[string]any{
				"host": "db-primary.example.com",
				"port": 5432,
				"ssl":  true,
			},
			"replica": map[string]any{
				"host": "db-replica.example.com",
				"port": 5432,
				"ssl":  false,
			},
		},
		"services": []map[string]any{
			{
				"name": "user-service",
				"port": 8081,
			},
			{
				"name": "order-service",
				"port": 8082,
			},
		},
	}

	source, err := NewValueSourceFromMap(complexInput)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	expected := map[string]any{
		"application.features.auth":  true,
		"application.features.cache": false,
		"application.name":           "my-app",
		"application.version":        "2.1.0",
		"databases.primary.host":     "db-primary.example.com",
		"databases.primary.port":     5432,
		"databases.primary.ssl":      true,
		"databases.replica.host":     "db-replica.example.com",
		"databases.replica.port":     5432,
		"databases.replica.ssl":      false,
		"services": []map[string]any{
			{
				"name": "user-service",
				"port": 8081,
			},
			{
				"name": "order-service",
				"port": 8082,
			},
		},
	}
	assert.Equal(t, expected, result)

	// 평면화된 키로 값에 접근 가능한지 확인
	assert.Equal(t, "my-app", result["application.name"])
	assert.Equal(t, "db-primary.example.com", result["databases.primary.host"])
	assert.Equal(t, true, result["application.features.auth"])
}

func TestValueSource_EdgeCases(t *testing.T) {
	t.Run("empty nested map", func(t *testing.T) {
		input := map[string]any{
			"empty": map[string]any{},
		}

		source, err := NewValueSourceFromMap(input)
		require.NoError(t, err)

		result, err := source.Load()
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})

	t.Run("mixed types", func(t *testing.T) {
		input := map[string]any{
			"string": "hello",
			"int":    42,
			"float":  3.14,
			"bool":   true,
			"slice":  []string{"a", "b", "c"},
		}

		source, err := NewValueSourceFromMap(input)
		require.NoError(t, err)

		result, err := source.Load()
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})

	t.Run("nil values in map", func(t *testing.T) {
		input := map[string]any{
			"nil_value": nil,
			"valid":     "value",
		}

		source, err := NewValueSourceFromMap(input)
		require.NoError(t, err)

		result, err := source.Load()
		require.NoError(t, err)
		assert.Equal(t, input, result)
	})
}
