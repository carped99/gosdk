package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempFile(t *testing.T, content string, ext string) string {
	t.Helper()
	tmp, err := os.CreateTemp("", "testfile-*"+ext)
	require.NoError(t, err)
	_, err = tmp.Write([]byte(content))
	require.NoError(t, err)
	require.NoError(t, tmp.Close())
	return tmp.Name()
}

func TestFileSource_Load_JSON(t *testing.T) {
	jsonContent := `{"name": "myapp", "port": 8080, "debug": true}`
	file := writeTempFile(t, jsonContent, ".json")
	defer os.Remove(file)

	source, err := NewFileSource(WithFilePath(file))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"name":  "myapp",
		"port":  float64(8080),
		"debug": true,
	}
	assert.Equal(t, expected, result)
}

func TestFileSource_Load_YAML(t *testing.T) {
	yamlContent := `name: myapp
port: 8080
debug: true`
	file := writeTempFile(t, yamlContent, ".yaml")
	defer os.Remove(file)

	source, err := NewFileSource(WithFilePath(file))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"name":  "myapp",
		"port":  8080,
		"debug": true,
	}
	assert.Equal(t, expected, result)
}

func TestFileSource_Load_MultipleFiles(t *testing.T) {
	jsonContent := `{"foo": "bar"}`
	yamlContent := `baz: qux`
	file1 := writeTempFile(t, jsonContent, ".json")
	file2 := writeTempFile(t, yamlContent, ".yaml")
	defer os.Remove(file1)
	defer os.Remove(file2)

	source, err := NewFileSource(WithFilePath(file1, file2))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"foo": "bar",
		"baz": "qux",
	}
	assert.Equal(t, expected, result)
}

func TestFileSource_Load_UnsupportedExtension(t *testing.T) {
	txtContent := `hello: world`
	file := writeTempFile(t, txtContent, ".txt")
	defer os.Remove(file)

	source, err := NewFileSource(WithFilePath(file))
	require.NoError(t, err)

	_, err = source.Load()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported file extension")
}

func TestFileSource_Load_FileNotExist(t *testing.T) {
	source, err := NewFileSource(WithFilePath("not_exist.json"))
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{}, result)
}

func TestFileSource_Load_EmptyFiles(t *testing.T) {
	source, err := NewFileSource()
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{}, result)
}

func TestFileSource_Watch(t *testing.T) {
	source, err := NewFileSource()
	require.NoError(t, err)

	watcher, err := source.Watch()
	require.NoError(t, err)
	assert.Nil(t, watcher)
}

func TestFileSource_Option_Error(t *testing.T) {
	errOpt := func(fs *FileSource) error {
		return assert.AnError
	}
	_, err := NewFileSource(errOpt)
	assert.Error(t, err)
}
