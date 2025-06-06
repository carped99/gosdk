package events

import (
	"encoding/json"
	"entgo.io/ent"
	"fmt"
	"github.com/huandu/xstrings"
	"reflect"
)

// NewSimpleEvent creates a new message for domain events based on the operation type
func NewSimpleEvent(eventDomain, eventType string, value any) (*Message, error) {
	objectType := getObjectTypeFromValue(value)

	payload, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}

	return NewMessageBuilder().
		SetEventTopic(eventDomain + ".events").
		SetEventDomain(eventDomain).
		SetEventType(eventType).
		SetObjectType(objectType).
		SetPayload(payload).
		Build()
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
