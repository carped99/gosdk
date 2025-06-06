package entx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type NameDescriptionMixin struct {
	mixin.Schema
}

func NameField() ent.Field {
	return field.
		String("name").
		NotEmpty().
		Annotations(
			entgql.OrderField("NAME"),
		)
}

func DescriptionField() ent.Field {
	return field.
		String("description").
		Optional().
		Nillable().
		Annotations(
			entgql.OrderField("DESCRIPTION"),
		)
}

func EnabledField() ent.Field {
	return field.
		Bool("enabled").
		Default(true)
}

func AuditFields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").Immutable().Annotations(
			entgql.Type("DateTime"),
			entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
		),

		field.UUID("created_by", uuid.UUID{}).Optional().Nillable().Annotations(
			entgql.Type("ID"),
			entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
		),

		field.Time("updated_at").Annotations(
			entgql.Type("DateTime"),
			entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
		),

		field.UUID("updated_by", uuid.UUID{}).Optional().Nillable().Annotations(
			entgql.Type("ID"),
			entgql.Skip(entgql.SkipMutationCreateInput, entgql.SkipMutationUpdateInput),
		),
	}
}

func (NameDescriptionMixin) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			Annotations(entgql.OrderField("NAME")),
		field.String("description").
			Optional().
			Nillable().
			Annotations(entgql.OrderField("DESCRIPTION")),
	}
}
