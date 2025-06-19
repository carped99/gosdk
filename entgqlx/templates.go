package entgqlx

import (
	"embed"
	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
)

var (
	//go:embed template/*
	_templates embed.FS

	FederationTemplate  = parseT("template/federation.tmpl")
	FilterInputTemplate = parseT("template/filter_input.tmpl")

	//TemplateFuncs = template.FuncMap{
	//	"fieldCollections":    fieldCollections,
	//	"fieldMapping":        fieldMapping,
	//	"filterEdges":         filterEdges,
	//	"filterFields":        filterFields,
	//	"filterNodes":         filterNodes,
	//	"gqlIDType":           gqlIDType,
	//	"gqlMarshaler":        gqlMarshaler,
	//	"gqlUnmarshaler":      gqlUnmarshaler,
	//	"hasWhereInput":       hasWhereInput,
	//	"isRelayConn":         isRelayConn,
	//	"isSkipMode":          isSkipMode,
	//	"mutationInputs":      mutationInputs,
	//	"nodeImplementors":    nodeImplementors,
	//	"nodeImplementorsVar": nodeImplementorsVar,
	//	"nodePaginationNames": nodePaginationNames,
	//	"orderFields":         orderFields,
	//	"skipMode":            skipModeFromString,
	//}
)

func parseT(path string) *gen.Template {
	return gen.MustParse(gen.NewTemplate(path).
		Funcs(entgql.TemplateFuncs).
		ParseFS(_templates, path))
}
