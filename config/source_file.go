package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/v2"
)

type FileSourceOption func(*FileSource) error

func WithFilePath(files ...string) FileSourceOption {
	return func(s *FileSource) error {
		for _, filePath := range files {
			if filePath != "" {
				s.files = append(s.files, filePath)
			}
		}
		return nil
	}
}

// FileSource provides file-based configuration
type FileSource struct {
	files []string
}

// NewFileSource creates a new FileSource with the given files
func NewFileSource(opts ...FileSourceOption) (*FileSource, error) {
	fs := &FileSource{
		files: make([]string, 0),
	}

	// Apply options to the FileSource
	for _, opt := range opts {
		if err := opt(fs); err != nil {
			return nil, err
		}
	}

	return fs, nil
}

// Load returns the configuration data from files
func (fs *FileSource) Load() (map[string]any, error) {
	result := make(map[string]any)

	if len(fs.files) == 0 {
		return result, nil
	}

	for _, filePath := range fs.files {
		if !fs.fileExists(filePath) {
			continue // 파일이 없으면 건너뛰기
		}

		data, err := fs.loadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("파일 로드 실패 (%s): %w", filePath, err)
		}

		// 파일 데이터를 결과에 병합
		for k, v := range data {
			result[k] = v
		}
	}

	return result, nil
}

// loadFile loads a single file and returns its data as map
func (fs *FileSource) loadFile(filePath string) (map[string]any, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	// koanf 인스턴스 생성
	k := koanf.New(".")

	// 확장된 파일 newEnvFileProvider 사용
	provider := newEnvFileProvider(filePath)

	// 파일 확장자에 따라 파서 선택하고 로드
	switch ext {
	case ".json":
		err := k.Load(provider, json.Parser())
		if err != nil {
			return nil, err
		}
	case ".yaml", ".yml":
		err := k.Load(provider, yaml.Parser())
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", ext)
	}

	return k.All(), nil
}

// fileExists checks if a file exists
func (fs *FileSource) fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

func (fs *FileSource) Watch() (Watcher, error) {
	// no-op
	return nil, nil
}
