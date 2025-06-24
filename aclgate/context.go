package aclgate

import (
	"context"
	"fmt"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const (
	// serviceKey is the key used to store ClientService in context
	serviceKey contextKey = "_acl_client_service_"
)

// FromContext retrieves ClientService from the context
func FromContext(ctx context.Context) (ClientService, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	service, ok := ctx.Value(serviceKey).(ClientService)
	if !ok {
		return nil, fmt.Errorf("ClientService not found in context or invalid type")
	}
	return service, nil
}

// MustFromContext retrieves ClientService from the context and panics if not found
func MustFromContext(ctx context.Context) ClientService {
	service, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}
	return service
}

// NewContext creates a new context with ClientService
func NewContext(ctx context.Context, service ClientService) context.Context {
	return context.WithValue(ctx, serviceKey, service)
}

// Check is a helper function to check permissions using the service from context
func Check(ctx context.Context, tuple Tuple) (bool, error) {
	return checkTuple(MustFromContext(ctx), ctx, tuple)
}

// BatchCheck is a helper function to check multiple permissions using the service from context
func BatchCheck(ctx context.Context, reqs []CheckRequest) ([]BatchCheckResult, error) {
	service := MustFromContext(ctx)
	return service.BatchCheck(ctx, reqs)
}

// Mutate is a helper function to mutate permissions using the service from context
func Mutate(ctx context.Context, writes, deletes []Tuple) (bool, error) {
	service := MustFromContext(ctx)
	return service.Mutate(ctx, writes, deletes)
}

func CanCreate(ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return canCreate(MustFromContext(ctx), ctx, resource, subject)
}

func CanRead(ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return canRead(MustFromContext(ctx), ctx, resource, subject)
}

func CanWrite(ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return canWrite(MustFromContext(ctx), ctx, resource, subject)
}

func CanDelete(ctx context.Context, resource Resource, subject Subject) (bool, error) {
	return canDelete(MustFromContext(ctx), ctx, resource, subject)
}

// Write is a helper function to write permissions using the service from context
func Write(ctx context.Context, tuples []Tuple) (bool, error) {
	return writeTuples(MustFromContext(ctx), ctx, tuples)
}

// Delete is a helper function to delete permissions using the service from context
func Delete(ctx context.Context, tuples []Tuple) (bool, error) {
	return deleteTuples(MustFromContext(ctx), ctx, tuples)
}

// DeleteResource is a helper function to delete all permissions for a resource using the service from context
func DeleteResource(ctx context.Context, resourceType, resourceId string) (bool, error) {
	return deleteResource(MustFromContext(ctx), ctx, NewResource(resourceType, resourceId))
}

// DeleteSubject is a helper function to delete all permissions for a subject using the service from context
func DeleteSubject(ctx context.Context, subjectType, subjectId string) (bool, error) {
	return deleteSubject(MustFromContext(ctx), ctx, NewSubject(subjectType, subjectId))
}
