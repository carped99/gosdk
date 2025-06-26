package loader

import (
	"fmt"
	"github.com/carped99/gosdk/config"
	"github.com/carped99/gosdk/config/source"
	"github.com/spf13/pflag"
	"log/slog"
	"reflect"
)

type LoadOption func(*loadOptions)

type loadOptions struct {
	envOpts  []source.EnvSourceOption
	fileOpts []source.FileSourceOption
	defaults any
}

// newLoadConfig 설정 파일을 읽기 위한 Config 객체를 생성
func newLoadConfig(flags *pflag.FlagSet, opts ...LoadOption) (config.Config, error) {
	options := &loadOptions{
		fileOpts: []source.FileSourceOption{
			source.WithFilePath("config.yaml"),
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	var sources []config.Source

	if options.defaults != nil {
		if src, err := source.NewValueSourceFromStruct(options.defaults, "koanf"); err != nil {
			return nil, err
		} else {
			sources = append(sources, src)
		}
	}

	if src, err := source.NewEnvSource(options.envOpts...); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	if src, err := source.NewFileSource(options.fileOpts...); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	if src, err := source.NewFlagSource(flags); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	return config.NewConfig(
		config.WithSource(sources...),
	)
}

func WithDefaultValueOptions(defaults any) LoadOption {
	return func(o *loadOptions) {
		o.defaults = defaults
	}
}

func WithEnvOptions(opts ...source.EnvSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.envOpts = opts
	}
}

func WithFileOptions(opts ...source.FileSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.fileOpts = opts
	}
}

func WithFlagOptions(opts ...source.FileSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.fileOpts = opts
	}
}

func LoadConfig[T any](defaults T, flags *pflag.FlagSet, opts ...LoadOption) (*T, error) {
	var result T
	typ := reflect.TypeOf(result)
	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("LoadConfig: T must be a struct type, got %s", typ.Kind())
	}

	// 기본값 옵션 설정
	opts = append(opts, WithDefaultValueOptions(defaults))

	c, err := newLoadConfig(flags, opts...)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := c.Close(); err != nil {
			slog.Error("failed to close config", slog.Any("error", err))
		}
	}()

	if err := c.Load(); err != nil {
		return nil, err
	}

	if err := c.Scan(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
