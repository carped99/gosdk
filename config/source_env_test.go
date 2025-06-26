package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEnvSource_Load_NoPrefix(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("APP_NAME", "myapp")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("APP_DEBUG", "true")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")
		os.Unsetenv("APP_DEBUG")
	}()

	source, err := NewEnvSource()
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	assert.Contains(t, result, "app.name")
	assert.Contains(t, result, "app.port")
	assert.Contains(t, result, "app.debug")
	assert.Equal(t, "myapp", result["app.name"])
	assert.Equal(t, "8080", result["app.port"])
	assert.Equal(t, "true", result["app.debug"])
}

func TestEnvSource_Load_WithPrefix(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("TEST_APP_NAME", "myapp")
	os.Setenv("TEST_APP_PORT", "8080")
	os.Setenv("TEST_APP_DEBUG", "true")
	defer func() {
		os.Unsetenv("TEST_APP_NAME")
		os.Unsetenv("TEST_APP_PORT")
		os.Unsetenv("TEST_APP_DEBUG")
	}()

	source, err := NewEnvSource(WithEnvPrefix("TEST"))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"app.name":  "myapp",
		"app.port":  "8080",
		"app.debug": "true",
	}
	assert.Equal(t, expected, result)
}

func TestEnvSource_Load_WithOverrides(t *testing.T) {
	// 환경 변수 설정
	os.Setenv("TEST_APP_NAME", "myapp")
	os.Setenv("TEST_APP_PORT", "8080")
	defer func() {
		os.Unsetenv("TEST_APP_NAME")
		os.Unsetenv("TEST_APP_PORT")
	}()

	nameMappings := map[string]string{
		"APP_NAME": "override_name",
		"APP_PORT": "override_port",
	}

	source, err := NewEnvSource(
		WithEnvPrefix("TEST"),
		WithEnvNameMapping(nameMappings),
	)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"override.name": "myapp",
		"override.port": "8080",
	}
	assert.Equal(t, expected, result)
}

func TestEnvSource_Load_Empty(t *testing.T) {
	source, err := NewEnvSource(WithEnvPrefix("NONEXISTENT"))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{}, result)
}

func TestEnvSource_Watch(t *testing.T) {
	source, err := NewEnvSource()
	require.NoError(t, err)

	watcher, err := source.Watch()
	require.NoError(t, err)
	assert.Nil(t, watcher)
}

func TestEnvSource_Option_Error(t *testing.T) {
	errOpt := func(es *EnvSource) error {
		return assert.AnError
	}
	_, err := NewEnvSource(errOpt)
	assert.Error(t, err)
}

func TestEnvSource_KeyTransformation(t *testing.T) {
	// 언더스코어가 점으로 변환되는지 테스트
	os.Setenv("TEST_DATABASE_HOST", "localhost")
	os.Setenv("TEST_DATABASE_PORT", "5432")
	os.Setenv("TEST_REDIS_HOST", "127.0.0.1")
	defer func() {
		os.Unsetenv("TEST_DATABASE_HOST")
		os.Unsetenv("TEST_DATABASE_PORT")
		os.Unsetenv("TEST_REDIS_HOST")
	}()

	source, err := NewEnvSource(WithEnvPrefix("TEST"))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"database.host": "localhost",
		"database.port": "5432",
		"redis.host":    "127.0.0.1",
	}
	assert.Equal(t, expected, result)
}

func TestEnvSource_ComplexKeys(t *testing.T) {
	// 복잡한 키 구조 테스트
	os.Setenv("TEST_APP_DATABASE_PRIMARY_HOST", "db1.example.com")
	os.Setenv("TEST_APP_DATABASE_PRIMARY_PORT", "5432")
	os.Setenv("TEST_APP_DATABASE_REPLICA_HOST", "db2.example.com")
	defer func() {
		os.Unsetenv("TEST_APP_DATABASE_PRIMARY_HOST")
		os.Unsetenv("TEST_APP_DATABASE_PRIMARY_PORT")
		os.Unsetenv("TEST_APP_DATABASE_REPLICA_HOST")
	}()

	source, err := NewEnvSource(WithEnvPrefix("TEST"))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"app.database.primary.host": "db1.example.com",
		"app.database.primary.port": "5432",
		"app.database.replica.host": "db2.example.com",
	}
	assert.Equal(t, expected, result)
}

func TestEnvSource_NoPrefix_WithSpecificVars(t *testing.T) {
	// prefix 없이 특정 환경 변수만 테스트
	os.Setenv("APP_NAME", "myapp")
	os.Setenv("APP_PORT", "8080")
	defer func() {
		os.Unsetenv("APP_NAME")
		os.Unsetenv("APP_PORT")
	}()

	source, err := NewEnvSource()
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	// 시스템 환경 변수도 포함되므로 우리가 설정한 변수들이 포함되어 있는지 확인
	assert.Contains(t, result, "app.name")
	assert.Contains(t, result, "app.port")
	assert.Equal(t, "myapp", result["app.name"])
	assert.Equal(t, "8080", result["app.port"])
}
