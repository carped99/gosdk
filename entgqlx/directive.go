package entgqlx

import (
	"entgo.io/contrib/entgql"
	"github.com/vektah/gqlparser/v2/ast"
)

func KeyDirective(fields string) entgql.Directive {
	var args []*ast.Argument
	if fields != "" {
		args = append(args, &ast.Argument{
			Name: "fields",
			Value: &ast.Value{
				Raw:  fields,
				Kind: ast.StringValue,
			},
		})
	}

	return entgql.NewDirective("key", args...)
}
