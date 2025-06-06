package bootstrap

type LogConfig struct {
	Level            string                 `json:"level"`              // log level: debug, info, warn, error
	Format           string                 `json:"format"`             // json or text
	Caller           bool                   `json:"caller"`             // 출력 시 caller 정보 포함 여부
	TimeFormat       string                 `json:"time_format"`        // 로그 타임스탬프 형식
	Stacktrace       bool                   `json:"stacktrace"`         // 스택트레이스 포함 여부
	StacktraceLevel  string                 `json:"stacktrace_level"`   // 스택트레이스를 포함할 최소 로그 레벨
	Development      bool                   `json:"development"`        // 개발 모드 여부
	InitialFields    map[string]interface{} `json:"initial_fields"`     // 초기 로그 필드
	OutputPaths      []string               `json:"output_paths"`       // 로그 출력 경로
	ErrorOutputPaths []string               `json:"error_output_paths"` // 에러 로그 출력 경로
}

var DefaultLogConfig = LogConfig{
	Level:            "info",
	Format:           "json",
	Caller:           true,
	TimeFormat:       "2006-01-02T15:04:05.000Z07:00",
	Stacktrace:       true,
	StacktraceLevel:  "error",
	Development:      false,
	InitialFields:    map[string]interface{}{},
	OutputPaths:      []string{"stdout"},
	ErrorOutputPaths: []string{"stderr"},
}
