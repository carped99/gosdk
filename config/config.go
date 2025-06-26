package config

import (
	"fmt"
	"sync"

	"github.com/knadh/koanf/providers/confmap"
	"github.com/knadh/koanf/v2"
)

type Config interface {
	Load() error
	Get(key string) Value
	Exists(key string) bool
	Scan(v any) error
	ScanPath(path string, v any) error
	All() map[string]any
	Close() error
}

type config struct {
	mu        sync.RWMutex
	k         *koanf.Koanf
	delimiter string
	sources   []Source
	watchers  []Watcher
}

func (c *config) Close() error {
	for _, w := range c.watchers {
		if err := w.Stop(); err != nil {
			return err
		}
	}
	return nil
}

type Option func(*config) error

func NewConfig(opts ...Option) (Config, error) {
	c := &config{
		delimiter: ".",
		k:         koanf.New("."),
	}
	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("failed to apply config option: %w", err)
		}
	}
	return c, nil
}

func WithSource(sources ...Source) Option {
	return func(c *config) error {
		for _, it := range sources {
			if it == nil {
				return fmt.Errorf("source cannot be nil")
			}
			c.sources = append(c.sources, it)
		}
		return nil
	}
}

func WithDelimiter(delim string) Option {
	return func(c *config) error {
		if delim == "" {
			return fmt.Errorf("delimiter cannot be empty")
		}
		c.delimiter = delim
		c.k = koanf.New(delim)
		return nil
	}
}

// Get Value
func (c *config) Get(key string) Value {
	c.mu.RLock()
	defer c.mu.RUnlock()
	av := &atomicValue{}
	av.Store(c.k.Get(key))
	return av
}

// Load loads configuration from all sources
func (c *config) Load() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, src := range c.sources {
		data, err := src.Load()
		if err != nil {
			return fmt.Errorf("failed to load source %T: %w", src, err)
		}
		if len(data) > 0 {
			if err := c.k.Load(confmap.Provider(data, c.k.Delim()), nil); err != nil {
				return fmt.Errorf("failed to load configuration from source %T: %w", src, err)
			}
		}
	}
	return nil
}

func (c *config) Scan(v any) error {
	return c.ScanPath("", v)
}

func (c *config) ScanPath(path string, v any) error {
	if v == nil {
		return fmt.Errorf("value cannot be nil")
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.k.Unmarshal(path, v)
}

func (c *config) All() map[string]any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.k.All()
}

func (c *config) Exists(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.k.Exists(key)
}
