package config

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlagSource_Load(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("name", "default", "name of the app")
	fs.Int("port", 8080, "port number")
	fs.Bool("debug", false, "enable debug mode")

	_ = fs.Set("name", "myapp")
	_ = fs.Set("port", "1234")
	_ = fs.Set("debug", "true")

	source, err := NewFlagSource(fs)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"name":  "myapp",
		"port":  1234,
		"debug": true,
	}
	assert.Equal(t, expected, result)
}

func TestFlagSource_Load_Defaults(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fs.String("foo", "bar", "default value")
	fs.Int("num", 42, "default int")

	source, err := NewFlagSource(fs)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)

	expected := map[string]any{
		"foo": "bar",
		"num": 42,
	}
	assert.Equal(t, expected, result)
}

func TestFlagSource_Load_NilFlagSet(t *testing.T) {
	source, err := NewFlagSource(nil)
	require.NoError(t, err)

	result, err := source.Load()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{}, result)
}

func TestFlagSource_Watch(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	source, err := NewFlagSource(fs)
	require.NoError(t, err)

	watcher, err := source.Watch()
	require.NoError(t, err)
	assert.Nil(t, watcher)
}

func TestFlagSource_Option(t *testing.T) {
	called := false
	opt := func(fs *FlagSource) error {
		called = true
		return nil
	}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	_, err := NewFlagSource(fs, opt)
	require.NoError(t, err)
	assert.True(t, called)
}

func TestFlagSource_Option_Error(t *testing.T) {
	errOpt := func(fs *FlagSource) error {
		return assert.AnError
	}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	_, err := NewFlagSource(fs, errOpt)
	assert.Error(t, err)
}
