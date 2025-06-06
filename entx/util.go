package entx

import "entgo.io/ent"

func JoinFields(values ...[]ent.Field) []ent.Field {
	var fields []ent.Field
	for _, fs := range values {
		fields = append(fields, fs...)
	}
	return fields
}
