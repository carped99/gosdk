package config

import (
	"fmt"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

// EnvSource provides environment variable configuration
type EnvSource struct {
	prefix       string
	nameMappings map[string]string
}

type EnvSourceOption func(*EnvSource) error

// NewEnvSource creates a new EnvSource with the given prefix
func NewEnvSource(opts ...EnvSourceOption) (*EnvSource, error) {
	es := &EnvSource{
		prefix:       "",
		nameMappings: make(map[string]string),
	}

	for _, opt := range opts {
		if err := opt(es); err != nil {
			return nil, err
		}
	}

	return es, nil
}

// WithEnvPrefix allows setting a prefix for environment variable names
func WithEnvPrefix(prefix string) EnvSourceOption {
	return func(es *EnvSource) error {
		es.prefix = prefix
		return nil
	}
}

// WithEnvNameMapping allows setting nameMappings for environment variable names
func WithEnvNameMapping(nameMappings map[string]string) EnvSourceOption {
	return func(es *EnvSource) error {
		es.nameMappings = nameMappings
		return nil
	}
}

// Load returns the configuration data from environment variables
func (es *EnvSource) Load() (map[string]any, error) {
	k := koanf.New(".")

	err := k.Load(env.Provider(es.prefix, k.Delim(), func(s string) string {
		if es.prefix != "" {
			s = strings.TrimPrefix(s, es.prefix+"_")
		}

		if v, ok := es.nameMappings[s]; ok {
			return strings.ReplaceAll(v, "_", k.Delim())
		}

		return strings.ReplaceAll(strings.ToLower(s), "_", k.Delim())
	}), nil)

	if err != nil {
		return nil, fmt.Errorf("failed to load env source: %w", err)
	}

	return k.All(), nil
}

func (es *EnvSource) Watch() (Watcher, error) {
	// no-op
	return nil, nil
}
