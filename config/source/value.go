package source

import (
	"fmt"
	"github.com/carped99/gosdk/config"
	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// ValueSource provides default configuration values
type ValueSource struct {
	value map[string]any
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
func NewValueSourceFromStruct(defaults any, tag string) (*ValueSource, error) {
	if defaults == nil {
		return &ValueSource{
			value: map[string]any{},
		}, nil
	}

	k := koanf.New(".")
	if err := k.Load(structs.Provider(defaults, tag), nil); err != nil {
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

func (ds *ValueSource) Watch() (config.Watcher, error) {
	// no-op
	return nil, nil
}
