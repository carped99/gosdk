package config

import (
	"fmt"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

type ValueSourceOption func(*ValueSource) error

// ValueSource provides default configuration values
type ValueSource struct {
	value map[string]any
	tag   string
}

// NewValueSourceFromMap creates a new ValueSource with the given value
func NewValueSourceFromMap(value map[string]any) (*ValueSource, error) {
	k := koanf.New(".")
	if err := k.Load(confmap.Provider(value, k.Delim()), nil); err != nil {
		return nil, fmt.Errorf("failed to load map source: %w", err)
	}
	return &ValueSource{
		value: k.All(),
	}, nil
}

// NewValueSourceFromStruct creates a ValueSource from a struct
func NewValueSourceFromStruct(defaults any, opts ...ValueSourceOption) (*ValueSource, error) {
	vs := &ValueSource{
		value: map[string]any{},
		tag:   "json",
	}

	if defaults == nil {
		return vs, nil
	}

	for _, opt := range opts {
		if err := opt(vs); err != nil {
			return nil, err
		}
	}

	k := koanf.New(".")
	if err := k.Load(structs.Provider(defaults, vs.tag), nil); err != nil {
		return nil, fmt.Errorf("failed to load struct source: %w", err)
	}

	return &ValueSource{
		value: k.All(),
	}, nil
}

// Load returns the default configuration data
func (ds *ValueSource) Load() (map[string]any, error) {
	return ds.value, nil
}

func (ds *ValueSource) Watch() (Watcher, error) {
	// no-op
	return nil, nil
}

// WithValueTag sets the tag used for struct fields in ValueSource
func WithValueTag(tag string) ValueSourceOption {
	return func(vs *ValueSource) error {
		vs.tag = tag
		return nil
	}
}
