package bootstrap

import "time"

type ServerConfig struct {
	HTTP  HTTPConfig       `json:"http" mapstructure:"http" env:"http"`
	HTTPS HTTPSConfig      `json:"https" mapstructure:"https" env:"https"`
	GRPC  GRPCServerConfig `json:"grpc" mapstructure:"grpc" env:"grpc"`
}

type HTTPConfig struct {
	Enabled           bool          `json:"enabled" mapstructure:"enabled" env:"enabled"` // TLS 사용 여부
	Address           string        `json:"address" mapstructure:"address" env:"address"` // 예: ":8080"
	ReadTimeout       time.Duration `json:"timeout.read" mapstructure:"timeout.read" env:"timeout.read"`
	ReadHeaderTimeout time.Duration `json:"timeout.header" mapstructure:"timeout.header" env:"timeout.header"`
	WriteTimeout      time.Duration `json:"timeout.write" mapstructure:"timeout.write" env:"timeout.write"`
	IdleTimeout       time.Duration `json:"timeout.idle" mapstructure:"timeout.idle" env:"timeout.idle"`
	ShutdownTimeout   time.Duration `json:"timeout.shutdown" mapstructure:"shutdown" env:"shutdown"`
	MaxHeaderBytes    int           `json:"max_header_bytes" mapstructure:"max_header_bytes" env:"max_header_bytes"`
}

type HTTPSConfig struct {
	Enabled           bool          `json:"enabled" mapstructure:"enabled" env:"enabled"`                                     // TLS 사용 여부
	Address           string        `json:"address" mapstructure:"address" env:"address"`                                     // 예: ":8443"
	CertFile          string        `json:"cert_file" mapstructure:"cert_file" env:"cert_file"`                               // PEM 형식
	KeyFile           string        `json:"key_file" mapstructure:"key_file" env:"key_file"`                                  // PEM 형식
	ClientCAs         string        `json:"client_ca_file" mapstructure:"client_ca_file" env:"client_ca_file"`                // mTLS 지원 시 필요
	RequireClientCert bool          `json:"require_client_cert" mapstructure:"require_client_cert" env:"require_client_cert"` // mTLS 여부
	ReadTimeout       time.Duration `json:"read_timeout" mapstructure:"read_timeout" env:"read_timeout"`
	ReadHeaderTimeout time.Duration `json:"header_timeout" mapstructure:"header_timeout" env:"header_timeout"`
	WriteTimeout      time.Duration `json:"write_timeout" mapstructure:"write_timeout" env:"write_timeout"`
	IdleTimeout       time.Duration `json:"idle_timeout" mapstructure:"idle_timeout" env:"idle_timeout"`
	MaxHeaderBytes    int           `json:"max_header_bytes" mapstructure:"max_header_bytes" env:"max_header_bytes"`
	ShutdownTimeout   time.Duration `json:"shutdown_timeout" mapstructure:"shutdown_timeout" env:"shutdown_timeout"`
}

type TLSConfig struct {
	CertFile string `json:"cert_file" mapstructure:"cert_file" env:"cert_file"`
	KeyFile  string `json:"key_file" mapstructure:"key_file" env:"key_file"`
	CAFile   string `json:"ca_file" mapstructure:"ca_file" env:"ca_file"`
}

var DefaultServerConfig = ServerConfig{
	HTTP: HTTPConfig{
		Enabled:           true,
		Address:           ":8080",
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
		MaxHeaderBytes:    1 << 20, // 1 MB
	},
	HTTPS: HTTPSConfig{
		Enabled:           false,
		Address:           ":8443",
		CertFile:          "",
		KeyFile:           "",
		ClientCAs:         "",
		RequireClientCert: false,
		ReadTimeout:       60 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       0,
		ReadHeaderTimeout: 0,
		MaxHeaderBytes:    1 << 20, // 1 MB
	},
	GRPC: DefaultGRPCServerConfig,
}
