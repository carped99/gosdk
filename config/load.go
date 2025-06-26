package config

import (
	"fmt"
	"github.com/spf13/pflag"
	"log/slog"
	"reflect"
)

type LoadOption func(*loadOptions)

type loadOptions struct {
	envOpts   []EnvSourceOption
	fileOpts  []FileSourceOption
	flagOpts  []FlagSourceOption
	valueOpts []ValueSourceOption
	defaults  any
}

// newLoadConfig 설정 파일을 읽기 위한 Config 객체를 생성
func newLoadConfig(flags *pflag.FlagSet, opts ...LoadOption) (Config, error) {
	options := &loadOptions{
		fileOpts: []FileSourceOption{
			WithFilePath("config.yaml"),
		},
		valueOpts: []ValueSourceOption{
			WithValueTag("json"),
		},
	}
	for _, opt := range opts {
		opt(options)
	}

	var sources []Source

	if options.defaults != nil {
		if src, err := NewValueSourceFromStruct(options.defaults, options.valueOpts...); err != nil {
			return nil, err
		} else {
			sources = append(sources, src)
		}
	}

	if src, err := NewEnvSource(options.envOpts...); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	if src, err := NewFileSource(options.fileOpts...); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	if src, err := NewFlagSource(flags, options.flagOpts...); err != nil {
		return nil, err
	} else {
		sources = append(sources, src)
	}

	return NewConfig(
		WithSource(sources...),
	)
}

func WithDefaultValueOptions(defaults any) LoadOption {
	return func(o *loadOptions) {
		o.defaults = defaults
	}
}

func WithValueTagOptions(opts ...ValueSourceOption) LoadOption {
	return func(o *loadOptions) {
		o.valueOpts = opts
	}
}

func WithEnvOptions(opts ...EnvSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.envOpts = opts
	}
}

func WithFileOptions(opts ...FileSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.fileOpts = opts
	}
}

func WithFlagOptions(opts ...FlagSourceOption) LoadOption {
	return func(c *loadOptions) {
		c.flagOpts = opts
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
