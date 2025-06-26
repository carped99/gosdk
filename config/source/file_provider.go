package source

import (
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"os"
	"regexp"
)

// envFileProvider 환경변수를 치환하는 파일 프로바이더
type envFileProvider struct {
	*file.File
}

// newEnvFileProvider creates a new envFileProvider
func newEnvFileProvider(filePath string) koanf.Provider {
	return &envFileProvider{
		File: file.Provider(filePath),
	}
}

// ReadBytes reads the contents of a file on disk and returns the bytes.
func (p *envFileProvider) ReadBytes() ([]byte, error) {
	originalContent, err := p.File.ReadBytes()
	if err != nil {
		return nil, err
	}

	// 환경변수 치환 적용
	expandedContent := p.expandEnvironmentVariables(string(originalContent))

	return []byte(expandedContent), nil
}

// expandEnvironmentVariables applies environment variable expansion to the content
func (p *envFileProvider) expandEnvironmentVariables(content string) string {
	// 기본값이 있는 환경변수 먼저 처리 (${VAR:default})
	defaultValueRegex := regexp.MustCompile(`\$\{([^:}]+):([^}]*)}`)
	result := defaultValueRegex.ReplaceAllStringFunc(content, func(match string) string {
		matches := defaultValueRegex.FindStringSubmatch(match)
		if len(matches) != 3 {
			return match // 매치되지 않으면 원본 반환
		}

		envVar := matches[1]
		defaultValue := matches[2]

		// 환경변수 값 가져오기
		envValue := os.Getenv(envVar)
		if envValue == "" {
			return defaultValue
		}
		return envValue
	})

	// 일반 환경변수 처리 (${VAR}) - os.ExpandEnv 사용
	result = os.ExpandEnv(result)

	return result
}
