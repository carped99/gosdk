package entx

import (
	"context"
	"entgo.io/ent"
	"time"
)

func WithTimestamps() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			now := time.Now()
			op := m.Op()
			switch op {
			case ent.OpCreate:
				_ = m.SetField("created_at", now)
				_ = m.SetField("updated_at", now)
			case ent.OpUpdate, ent.OpUpdateOne:
				_ = m.SetField("updated_at", now)
			case ent.OpDelete, ent.OpDeleteOne:
				_ = m.SetField("deleted_at", now)
			}
			return next.Mutate(ctx, m)
		})
	}
}
