package bootstrap

import "fmt"

type OpenFGAConfig struct {
	GRPC  GRPCClientConfig `json:"grpc"`
	Store StoreConfig      `json:"store"`
	Model ModelConfig      `json:"model"`
}

func (c *OpenFGAConfig) Validate() error {
	if c.Store.Id == "" {
		return fmt.Errorf("store ID must be provided")
	}
	return nil
}

type ApiConfig struct {
	Type string `json:"type"`
	Uri  string `json:"uri"`
}

type StoreConfig struct {
	Id string `json:"id"`
}

type ModelConfig struct {
	Id string `json:"id"`
}

var DefaultClientConfig = OpenFGAConfig{
	GRPC: DefaultGRPCClientConfig,
}
