package entx

import (
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	Config struct {
		enableAudit   bool
		enableTracing bool
	}

	Extension struct {
		entc.DefaultExtension
		config *Config
		hooks  []gen.Hook
	}

	ExtensionOption func(*Extension) error
)

func WithAudit() ExtensionOption {
	return func(c *Extension) error {
		c.config.enableAudit = true
		return nil
	}
}

func WithTracing() ExtensionOption {
	return func(c *Extension) error {
		c.config.enableTracing = true
		return nil
	}
}

func NewExtension(opts ...ExtensionOption) *Extension {
	ex := &Extension{
		config: &Config{},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			panic("failed to create extension: " + err.Error())
		}
	}

	//ex.hooks = append(ex.hooks, ex.genSchemaHook())
	return ex
}

//func (e *Extension) genSchemaHook() gen.Hook {
//	return func(next gen.Generator) gen.Generator {
//		return gen.GenerateFunc(func(g *gen.Graph) error {
//			schema, err := e.BuildSchema(g)
//			if err != nil {
//				return err
//			}
//
//			return os.WriteFile(e.path, []byte(printSchema(schema)), 0644)
//
//			return next.Generate(g)
//		})
//	}
//}

func (e *Extension) Options() []entc.Option {
	return []entc.Option{}
}

func (e *Extension) Templates() []*gen.Template {
	var templates []*gen.Template
	templates = append(templates, clientTemplate, defaultsTemplate, GetModelNameTemplate)
	templates = append(templates, WithTxTemplate)
	return templates
}

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func (e *Extension) Annotations() []entc.Annotation {
	return nil
}

var (
	_ entc.Extension = (*Extension)(nil)

	camel    = gen.Funcs["camel"].(func(string) string)
	pascal   = gen.Funcs["pascal"].(func(string) string)
	plural   = gen.Funcs["plural"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	snake    = gen.Funcs["snake"].(func(string) string)
)
