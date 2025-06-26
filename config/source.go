package config

// Source defines the interface for configuration sources
type Source interface {
	// Load returns the configuration data from this source
	Load() (map[string]any, error)
	Watch() (Watcher, error)
}

type Watcher interface {
	Next() (map[string]any, error)
	Stop() error
}
