package entx

import (
	"embed"
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
)

var (
	//go:embed template/*
	_templates embed.FS

	clientTemplate       = parseT("template/client.tmpl")
	defaultsTemplate     = parseT("template/defaults.tmpl")
	GetModelNameTemplate = parseT("template/get_model_name.tmpl")
	WithTxTemplate       = parseT("template/with_tx.tmpl")
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(entgql.TemplateFuncs).
		ParseFS(_templates, path))
}
