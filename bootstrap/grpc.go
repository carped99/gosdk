package bootstrap

import (
	"fmt"
	"time"
)

type GRPCClientConfig struct {
	Target         string                     `koanf:"target"`            // gRPC 서버 주소
	UseTLS         bool                       `koanf:"use_tls"`           // TLS 사용 여부
	TLS            *TLSConfig                 `koanf:"tls"`               // TLS 인증서 경로
	Retry          *GRPCClientRetryConfig     `koanf:"retry"`             // 재시도 설정
	KeepAlive      *GRPCClientKeepAliveConfig `koanf:"keepalive"`         // KeepAlive 설정
	MaxRecvMsgSize int                        `koanf:"max_recv_msg_size"` // 수신 메시지 최대 크기
	MaxSendMsgSize int                        `koanf:"max_send_msg_size"` // 송신 메시지 최대 크기
	DefaultTimeout time.Duration              `koanf:"default_timeout"`   // 기본 타임아웃
	Metadata       map[string]string          `koanf:"metadata"`          // 요청 시 메타데이터 헤더
}

func (c *GRPCClientConfig) Validate() error {
	if c.Target == "" {
		return fmt.Errorf("grpc target must be set")
	}
	return nil
}

type GRPCServerConfig struct {
	Enabled         bool          `json:"enabled" mapstructure:"enabled" env:"enabled"`
	Address         string        `json:"address" mapstructure:"address" env:"address"`
	CertFile        string        `json:"cert_file" mapstructure:"cert_file" env:"cert_file"`
	KeyFile         string        `json:"key_file" mapstructure:"key_file" env:"key_file"`
	CAFile          string        `json:"ca_file" mapstructure:"ca_file" env:"ca_file"`
	TLS             *TLSConfig    `json:"tls" mapstructure:"tls" env:"tls"`
	Reflection      bool          `json:"reflection" mapstructure:"reflection" env:"reflection"`
	ShutdownTimeout time.Duration `json:"shutdown_timeout" mapstructure:"shutdown_timeout" env:"shutdown_timeout"`
}

type GRPCClientRetryConfig struct {
	MaxAttempts       int           `koanf:"max_attempts"`       // 최대 시도 횟수
	InitialBackoff    time.Duration `koanf:"initial_backoff"`    // 최초 백오프
	MaxBackoff        time.Duration `koanf:"max_backoff"`        // 최대 백오프
	BackoffMultiplier float64       `koanf:"backoff_multiplier"` // 백오프 증가 비율
}

type GRPCClientKeepAliveConfig struct {
	Time                time.Duration `koanf:"time"`                  // Ping 주기
	Timeout             time.Duration `koanf:"timeout"`               // 응답 대기시간
	PermitWithoutStream bool          `koanf:"permit_without_stream"` // 스트림 없이 ping 허용 여부
}

// DefaultGRPCClientConfig is the default config for gRPC client.
var DefaultGRPCClientConfig = GRPCClientConfig{
	Target:         "localhost:50051",
	UseTLS:         false,
	TLS:            nil,
	Retry:          nil,
	KeepAlive:      nil,
	MaxRecvMsgSize: 1024 * 1024 * 4, // 4MB
	MaxSendMsgSize: 1024 * 1024 * 4, // 4MB
	DefaultTimeout: 5 * time.Second,
	Metadata:       nil,
}

var DefaultGRPCServerConfig = GRPCServerConfig{
	Enabled:         true,
	Address:         ":50051",
	CertFile:        "",
	KeyFile:         "",
	CAFile:          "",
	Reflection:      true,
	ShutdownTimeout: 10 * time.Second,
}
