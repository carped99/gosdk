package aclgate

import (
	"context"
	"fmt"
	"time"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey struct {
	name string
}

var (
	// serviceKey is the key used to store ClientService in context
	// Using struct type instead of string to prevent key collisions
	serviceKey = contextKey{name: "acl_client_service"}
)

// ContextOptions provides configuration for context operations
type ContextOptions struct {
	Timeout time.Duration
}

// DefaultContextOptions returns default context options
func DefaultContextOptions() *ContextOptions {
	return &ContextOptions{
		Timeout: 30 * time.Second,
	}
}

// FromContext retrieves ClientService from the context
func FromContext(ctx context.Context) (ClientService, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	service, ok := ctx.Value(serviceKey).(ClientService)
	if !ok {
		return nil, ErrServiceNotFound
	}

	if service == nil {
		return nil, fmt.Errorf("ClientService is nil in context")
	}

	return service, nil
}

// FromContextWithTimeout retrieves ClientService from context with timeout
func FromContextWithTimeout(ctx context.Context, timeout time.Duration) (ClientService, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}

	// Create a new context with timeout
	ctxWithTimeout, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Wait for context to be done or service to be available
	select {
	case <-ctxWithTimeout.Done():
		return nil, fmt.Errorf("context timeout or cancelled: %w", ctxWithTimeout.Err())
	default:
		return FromContext(ctx)
	}
}

// MustFromContext retrieves ClientService from the context and panics if not found
func MustFromContext(ctx context.Context) ClientService {
	service, err := FromContext(ctx)
	if err != nil {
		panic(fmt.Sprintf("failed to get ClientService from context: %v", err))
	}
	return service
}

// NewContext creates a new context with ClientService
func NewContext(ctx context.Context, service ClientService) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if service == nil {
		return ctx
	}

	return context.WithValue(ctx, serviceKey, service)
}

// NewContextWithTimeout creates a new context with ClientService and timeout
func NewContextWithTimeout(ctx context.Context, service ClientService, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(NewContext(ctx, service), timeout)
}

// Check is a helper function to check permissions using the service from context
func Check(ctx context.Context, tuple *Tuple) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTuple(service, ctx, tuple)
}

// CheckWithTimeout is a helper function to check permissions with timeout
func CheckWithTimeout(ctx context.Context, tuple *Tuple, timeout time.Duration) (bool, error) {
	service, err := FromContextWithTimeout(ctx, timeout)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTuple(service, ctx, tuple)
}

// BatchCheck is a helper function to check multiple permissions using the service from context
func BatchCheck(ctx context.Context, reqs []*CheckRequest) ([]*BatchCheckResult, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.BatchCheck(ctx, reqs)
}

// Mutate is a helper function to mutate permissions using the service from context
func Mutate(ctx context.Context, writes, deletes []*Tuple) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.Mutate(ctx, writes, deletes)
}

// CanCreate is a helper function to check if a subject can create a resource
func CanCreate(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTupleRelation(service, ctx, resource, subject, RelationCanCreate)
}

// CanRead is a helper function to check if a subject can read a resource
func CanRead(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTupleRelation(service, ctx, resource, subject, RelationCanRead)
}

// CanWrite is a helper function to check if a subject can write to a resource
func CanWrite(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTupleRelation(service, ctx, resource, subject, RelationCanWrite)
}

// CanDelete is a helper function to check if a subject can delete a resource
func CanDelete(ctx context.Context, resource *Resource, subject *Subject) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return checkTupleRelation(service, ctx, resource, subject, RelationCanDelete)
}

// Write is a helper function to write permissions using the service from context
func Write(ctx context.Context, tuples []*Tuple) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.Mutate(ctx, tuples, nil)
}

// Delete is a helper function to delete permissions using the service from context
func Delete(ctx context.Context, tuples []*Tuple) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.Mutate(ctx, nil, tuples)
}

// DeleteResource is a helper function to delete a resource using the service from context
func DeleteResource(ctx context.Context, resource *Resource) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.Mutate(ctx, nil, []*Tuple{
		{
			Resource: resource,
			Subject:  nil,
			Relation: nil,
		},
	})
}

// DeleteSubject is a helper function to delete a subject using the service from context
func DeleteSubject(ctx context.Context, subject *Subject) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to get service from context: %w", err)
	}
	return service.Mutate(ctx, nil, []*Tuple{
		{
			Resource: nil,
			Subject:  subject,
			Relation: nil,
		},
	})
}
