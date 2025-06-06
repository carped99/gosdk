package entx

import (
	"encoding/json"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

func OutboxFields() []ent.Field {
	return []ent.Field{
		field.Int64("id").Immutable(),
		field.UUID("uuid", uuid.UUID{}).Default(func() uuid.UUID {
			return uuid.New()
		}),
		field.String("event_topic").NotEmpty(),
		field.String("event_domain").NotEmpty(),
		field.String("object_type").NotEmpty(),
		field.String("event_type").NotEmpty(),
		field.String("producer").Optional().Nillable(),
		field.String("correlation_id").Optional().Nillable(),
		field.JSON("payload", json.RawMessage{}).SchemaType(map[string]string{
			"postgres": "jsonb",
		}),
		field.JSON("metadata", json.RawMessage{}).SchemaType(map[string]string{
			"postgres": "jsonb",
		}).Optional(),
		field.Time("created_at").Default(time.Now),
	}
}
