package entgqlx

import (
	"fmt"
	"reflect"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

type schemaGenerator struct {
	cfg            *config.Config
	path           string
	genSchema      bool
	genWhereFilter bool
}

var (
	inputObjectFilter = func(t string) bool {
		return strings.HasSuffix(t, "Input")
	}
)

type (
	FilterInputDefinition struct {
		Field      *gen.Field
		Definition *ast.Definition
	}
	FilterFieldDefinition struct {
		Field      *gen.Field
		Definition *ast.FieldDefinition
	}
)

func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb, formatter.WithIndent("  ")).
		FormatSchema(schema)
	return sb.String()
}

func (e *schemaGenerator) WhereFilter(g *gen.Graph, s *ast.Schema) error {
	for _, node := range g.Nodes {
		fmt.Println("Node: ", node.Name)
		if node.HasCompositeID() {
			continue
		}

		gqlType, ant, err := gqlTypeFromNode(node)
		if err != nil {
			return err
		}

		if !ant.Skip.Is(entgql.SkipWhereInput) {
			if err := e.buildFilterInput(node, gqlType, s); err != nil {
				return err
			}
		}
	}
	return nil
}

func (e *schemaGenerator) buildFilterInput(t *gen.Type, gqlType string, s *ast.Schema) error {
	// FilterInput 정의
	def, err := e.buildFilterInputDefinition(t)
	if err != nil {
		return err
	}

	if def == nil {
		return nil
	}

	s.AddTypes(def)

	fieldInputs, err := e.buildFieldFilterInputDefinition(t, gqlType)
	if err != nil {
		return err
	}

	for _, fieldFilter := range fieldInputs {
		s.AddTypes(fieldFilter.Definition)

		def.Fields = append(def.Fields, &ast.FieldDefinition{
			Name: camel(fieldFilter.Field.Name),
			Type: ast.NamedType(fieldFilter.Definition.Name, nil),
		})
	}

	return nil
}

func gqlTypeFromNode(t *gen.Type) (gqlType string, ant *entgql.Annotation, err error) {
	if ant, err = annotation(t.Annotations); err != nil {
		return
	}
	gqlType = t.Name
	if ant.Type != "" {
		gqlType = ant.Type
	}
	return
}

func annotation(ants gen.Annotations) (*entgql.Annotation, error) {
	ant := &entgql.Annotation{}
	if ants != nil && ants[ant.Name()] != nil {
		if err := ant.Decode(ants[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

func (e *schemaGenerator) buildFilterInputDefinition(t *gen.Type) (*ast.Definition, error) {
	name := t.Name + "FilterInput"
	def := &ast.Definition{
		Name: name,
		Kind: ast.InputObject,
	}

	// "not" 추가
	def.Fields = append(def.Fields, &ast.FieldDefinition{
		Name: "not",
		Type: ast.NamedType(name, nil),
	})

	// "and" and "or" 추가
	for _, op := range []string{"and", "or"} {
		def.Fields = append(def.Fields, &ast.FieldDefinition{
			Name: op,
			Type: ast.ListType(ast.NonNullNamedType(name, nil), nil),
		})
	}

	return def, nil
}

func (e *schemaGenerator) buildFieldFilterInputDefinition(t *gen.Type, nodeGQLType string) ([]*FilterInputDefinition, error) {
	definitions := make([]*FilterInputDefinition, 0)

	// 필드 목록
	fields := allFields(t)

	for _, f := range fields {
		if t.IsEdgeSchema() && f.IsEdgeField() || !f.Type.Comparable() || f.Sensitive() {
			continue
		}

		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(entgql.SkipWhereInput) {
			continue
		}

		name := pascal(t.Name + "_" + f.Name + "FilterInput")
		def := &ast.Definition{
			Name: name,
			Kind: ast.InputObject,
		}

		for _, op := range f.Ops() {
			fd := e.newFilterFieldDefinition(nodeGQLType, f, ant, op)
			def.Fields = append(def.Fields, fd.Definition)
		}

		definitions = append(definitions, &FilterInputDefinition{
			Field:      f,
			Definition: def,
		})
	}

	return definitions, nil
}

func (e *schemaGenerator) newFilterFieldDefinition(gqlType string, f *gen.Field, ant *entgql.Annotation, op gen.Op) *FilterFieldDefinition {
	name := camel(snake(op.Name()))
	def := &ast.FieldDefinition{
		Name: name,
	}

	//if e.scalarFunc != nil {
	//	if t := e.scalarFunc(f, op); t != "" {
	//		def.Type = namedType(t, true)
	//		return def
	//	}
	//}

	switch {
	case op.Niladic():
		def.Type = namedType("Boolean", true)
	case op.Variadic():
		def.Type = listNamedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	default:
		def.Type = namedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	}

	return &FilterFieldDefinition{
		Field:      f,
		Definition: def,
	}
}

func allFields(t *gen.Type) []*gen.Field {
	if t.ID == nil {
		return t.Fields
	}

	// NOTE: always keep the ID field at the beginning of the list
	return append([]*gen.Field{t.ID}, t.Fields...)
}

func (e *schemaGenerator) fieldDefinitionOp(gqlType string, f *gen.Field, ant *entgql.Annotation, op gen.Op) *ast.FieldDefinition {
	def := &ast.FieldDefinition{
		Name: camel(f.Name + "_" + op.Name()),
	}

	if op == gen.EQ {
		def.Name = camel(f.Name)
	}

	//if e.scalarFunc != nil {
	//	if t := e.scalarFunc(f, op); t != "" {
	//		def.Type = namedType(t, true)
	//		return def
	//	}
	//}

	switch {
	case op.Niladic():
		def.Type = namedType("Boolean", true)
	case op.Variadic():
		def.Type = listNamedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	default:
		def.Type = namedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	}
	return def
}

func namedType(name string, nullable bool) *ast.Type {
	if nullable {
		return ast.NamedType(name, nil)
	}
	return ast.NonNullNamedType(name, nil)
}

func listNamedType(name string, nullable bool) *ast.Type {
	t := ast.NonNullNamedType(name, nil)
	if nullable {
		return ast.ListType(t, nil)
	}
	return ast.NonNullListType(t, nil)
}

// mapScalar provides maps an ent.Schema type into GraphQL scalar type.
func (e *schemaGenerator) mapScalar(gqlType string, f *gen.Field, ant *entgql.Annotation, typeFilter func(string) bool) string {
	if ant != nil && ant.Type != "" {
		return ant.Type
	}
	scalar := f.Type.String()
	switch t := f.Type.Type; {
	case f.Name == "id":
		return "Id"
	case f.IsEdgeField():
		scalar = "Id"
	case t.Float():
		scalar = "Float"
	case t.Integer():
		scalar = "Int"
	case t == field.TypeString:
		scalar = "String"
	case t == field.TypeBool:
		scalar = "Boolean"
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		if typ, ok := e.hasMapping(f, typeFilter); ok {
			scalar = typ
		} else {
			scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		}
		if f.IsEnum() {
			// Use the GQL type as enum prefix. e.g. Todo.status
			// will generate an enum named "TodoStatus".
			scalar = gqlType + scalar
		}
		if f.Type.RType != nil && f.Type.RType.Name == "" {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				if strings.HasPrefix(f.Type.RType.Ident, "[]*") {
					scalar = "[" + scalar + "]"
				} else {
					scalar = "[" + scalar + "!]"
				}
			}
		}
	case t == field.TypeJSON:
		scalar = ""
		if f.Type.RType != nil {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				switch f.Type.RType.Ident {
				case "[]float64":
					scalar = "[Float!]"
				case "[]int":
					scalar = "[Int!]"
				case "[]string":
					scalar = "[String!]"
				}
			case reflect.Map:
				if f.Type.RType.Ident == "map[string]interface {}" {
					scalar = "Map"
				}
			}
		}
	}
	return scalar
}

// hasMapping reports if the gqlgen.yml has custom mapping for
// the given field type and returns its GraphQL name if exists.
func (e *schemaGenerator) hasMapping(f *gen.Field, typeFilter func(string) bool) (string, bool) {
	if e.cfg == nil {
		return "", false
	}

	// The string representation uses shortened package
	// names, and we override it for custom Go types.
	ident := f.Type.String()
	if idx := strings.IndexByte(ident, '.'); idx != -1 && f.HasGoType() && f.Type.PkgPath != "" {
		ident = f.Type.PkgPath + ident[idx:]
	}

	var gqlNames []string
	for t, v := range e.cfg.Models {
		for _, m := range v.Model {
			// A mapping was found from GraphQL name to field type.
			if strings.HasSuffix(m, ident) {
				gqlNames = append(gqlNames, t)
			}
		}
	}
	if count := len(gqlNames); count == 1 {
		return gqlNames[0], true
	} else if count > 1 && typeFilter != nil {
		for _, t := range gqlNames {
			if typeFilter(t) {
				return t, true
			}
		}
	}

	// If no custom mapping was found, fallback to the builtin scalar
	// types as mentioned in https://gqlgen.com/reference/scalars
	switch f.Type.String() {
	case "time.Time":
		return "Time", true
	case "map[string]interface{}":
		return "Map", true
	default:
		return "", false
	}
}
