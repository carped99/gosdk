package bootstrap

import "fmt"

type AclGateClientConfig struct {
	GRPC GRPCClientConfig `koanf:"grpc"`
}

func (c *AclGateClientConfig) Validate() error {
	if err := c.GRPC.Validate(); err != nil {
		return fmt.Errorf("aclgate config is invalid: %w", err)
	}
	return nil
}

// DefaultAclGateClientConfig is the default config for a client.
var DefaultAclGateClientConfig = AclGateClientConfig{
	GRPC: DefaultGRPCClientConfig,
}
