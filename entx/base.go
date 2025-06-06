package entx

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/mixin"
)

type BaseMixin struct {
	mixin.Schema
}

func (BaseMixin) Fields() []ent.Field {
	fields := []ent.Field{
		UUIDPrimary(),
	}

	fields = append(fields, NameDescriptionMixin{}.Fields()...)
	fields = append(fields, AuditMixin{}.Fields()...)

	return fields
}
