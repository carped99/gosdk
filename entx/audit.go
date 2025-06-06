package entx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
	"time"
)

type AuditMixin struct {
	mixin.Schema
}

func (AuditMixin) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now()).
			Immutable().
			Annotations(
				entgql.Type("DateTime"),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
		field.UUID("created_by", uuid.UUID{}).
			Optional().
			Nillable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
		field.Time("updated_at").Default(time.Now()).
			Annotations(
				entgql.Type("DateTime"),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
		field.UUID("updated_by", uuid.UUID{}).
			Optional().
			Nillable().
			Annotations(
				entgql.Type("ID"),
				entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
			),
	}
}
