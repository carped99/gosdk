package entx

import (
	"context"
	"entgo.io/ent"
	"fmt"
	"github.com/huandu/xstrings"
	"reflect"
	"runtime/debug"
	"shared/outbox"
)

func WithDomainEvent(eventDomain string, executor func(ctx context.Context) outbox.Executor) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (value ent.Value, err error) {
			value, err = next.Mutate(ctx, m)
			if err == nil {
				exec := executor(ctx)
				if exec == nil {
					return value, fmt.Errorf("executor not found in context")
				}

				var eventType string
				switch m.Op() {
				case ent.OpCreate:
					eventType = "created"
				case ent.OpUpdate, ent.OpUpdateOne:
					eventType = "updated"
				case ent.OpDelete, ent.OpDeleteOne:
					eventType = "deleted"
				}

				objectType := getObjectTypeFromValue(value)

				message, err := outbox.NewMessageBuilder().
					SetEventTopic(eventDomain + ".events").
					SetEventDomain(eventDomain).
					SetEventType(eventType).
					SetObjectType(objectType).
					SetPayload(value).
					Build()
				if err != nil {
					return value, err
				}

				err = outbox.NewPublisher(exec).Publish(ctx, message)
			}
			return value, err
		})
	}
}

func WithAclEvent(eventDomain string, executor func(ctx context.Context) outbox.Executor) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (value ent.Value, err error) {
			value, err = next.Mutate(ctx, m)
			if err == nil {
				exec := executor(ctx)
				if exec == nil {
					return value, fmt.Errorf("executor not found in context")
				}

				var eventType string
				switch m.Op() {
				case ent.OpCreate:
					eventType = "created"
				case ent.OpUpdate, ent.OpUpdateOne:
					eventType = "updated"
				case ent.OpDelete, ent.OpDeleteOne:
					eventType = "deleted"
				}

				objectType := getObjectTypeFromValue(value)

				tuple, err := outbox.NewAclTupleBuilder().
					SetResourceType(objectType).
					SetResourceID(getIdFromValue(value)).
					SetSubjectType("user").
					SetSubjectID("system").
					Build()
				if err != nil {
					return nil, err
				}

				payload, err := outbox.NewAclPayloadBuilder().
					AddTuples(tuple).
					Build()
				if err != nil {
					return nil, err
				}

				message, err := outbox.NewMessageBuilder().
					SetEventTopic("acls.events").
					SetEventDomain(eventDomain).
					SetEventType(eventType).
					SetObjectType(objectType).
					SetPayload(payload).
					Build()
				if err != nil {
					return value, err
				}

				err = outbox.NewPublisher(exec).Publish(ctx, message)
			}
			return value, err
		})
	}
}

func getModuleName() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	return info.Main.Path
}

// getObjectTypeFromValue extracts the object type from an ent.Value.
func getObjectTypeFromValue(v ent.Value) string {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return xstrings.ToSnakeCase(t.Name())
}

// getIdFromValue extracts the ID field from an ent.Value.
func getIdFromValue(v ent.Value) string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	idField := val.FieldByName("ID")
	if !idField.IsValid() {
		panic(fmt.Errorf("ID field not found in struct %T: struct must have an 'ID' field", v))
	}

	return fmt.Sprint(idField.Interface())
}
