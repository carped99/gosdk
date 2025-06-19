package entgqlx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

type (
	Config struct {
		enableAudit   bool
		enableTracing bool
	}

	Extension struct {
		schemaGenerator
		entc.DefaultExtension
		config      *Config
		hooks       []gen.Hook
		schemaHooks []entgql.SchemaHook
		templates   []*gen.Template
	}

	ExtensionOption func(*Extension) error
)

func WithFederation() ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = append(ex.templates, FederationTemplate)
		//ex.gqlSchemaHooks = append(ex.gqlSchemaHooks, removeNodeGoModel, removeNodeQueries, setPageInfoShareable)
		return nil
	}
}

func WithTemplates(templates ...*gen.Template) ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = templates
		return nil
	}
}

func WithSchemaGenerator() ExtensionOption {
	return func(e *Extension) error {
		e.genSchema = true
		return nil
	}
}

func WithSchemaPath(path string) ExtensionOption {
	return func(e *Extension) error {
		e.path = path
		return nil
	}
}

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
		schemaGenerator: schemaGenerator{
			genSchema: false,
		},
		config: &Config{}}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			panic("failed to apply extension option: " + err.Error())
		}
	}

	ex.schemaHooks = append(ex.schemaHooks, RemoveNodeGoModel, RemoveNodeQueries)
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
	return e.templates
}

func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

func (e *Extension) SchemaHooks() []entgql.SchemaHook {
	return e.schemaHooks
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
