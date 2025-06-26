package config

import (
	"fmt"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

type FlagSourceOption func(*FlagSource) error

// FlagSource provides command-line flag configuration
type FlagSource struct {
	flags *pflag.FlagSet
}

// NewFlagSource creates a new FlagSource with the given flag set
func NewFlagSource(flags *pflag.FlagSet, opts ...FlagSourceOption) (*FlagSource, error) {
	fs := &FlagSource{
		flags: flags,
	}

	for _, opt := range opts {
		if err := opt(fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

// Load returns the configuration data from command-line flags
func (fs *FlagSource) Load() (map[string]any, error) {
	if fs.flags == nil {
		return make(map[string]any), nil
	}

	k := koanf.New(".")
	err := k.Load(posflag.Provider(fs.flags, ".", k), nil)
	if err != nil {
		return nil, fmt.Errorf("플래그 로드 실패: %w", err)
	}

	return k.All(), nil
}

func (fs *FlagSource) Watch() (Watcher, error) {
	// no-op
	return nil, nil
}
