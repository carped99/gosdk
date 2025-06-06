package entx

import (
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

type UUIDMixin struct {
	mixin.Schema
}

func (UUIDMixin) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(func() uuid.UUID {
				return uuid.Must(uuid.NewV7())
			}).
			Unique().
			Immutable().
			Annotations(entgql.Type("ID"), entgql.OrderField("ID")),
	}
}

func UUIDPrimary(annotations ...schema.Annotation) ent.Field {
	builder := field.UUID("id", uuid.UUID{}).
		Default(func() uuid.UUID {
			return uuid.Must(uuid.NewV7()) // UUIDv7
		}).
		Unique().
		Immutable().
		Annotations(
			entgql.Type("ID"),       // GraphQL ID 타입
			entgql.OrderField("ID"), // 정렬 가능
		)

	if len(annotations) > 0 {
		builder = builder.Annotations(annotations...)
	}
	return builder
}

func UUIDForeign(name string, nillable bool, annotations ...schema.Annotation) ent.Field {
	builder := field.UUID(name, uuid.UUID{}).
		Annotations(entgql.Type("ID"))

	if nillable {
		builder = builder.Nillable().Optional()
	}

	if len(annotations) > 0 {
		builder = builder.Annotations(annotations...)
	}
	return builder
}

func UUIDField(name string, nillable bool, annotations ...schema.Annotation) ent.Field {
	builder := field.UUID(name, uuid.UUID{})

	if nillable {
		builder = builder.Nillable().Optional()
	}

	if len(annotations) > 0 {
		builder = builder.Annotations(annotations...)
	}
	return builder
}

func WithUUID() []ent.Field {
	return UUIDMixin{}.Fields()
}
