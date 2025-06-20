package aclgate

import (
	"context"
	"fmt"
	"net/http"
)

// contextKey is a private type for context keys to avoid collisions
type contextKey string

const (
	// aclServiceKey is the key used to store AclService in context
	aclServiceKey contextKey = "acl_service"
)

// WithContext creates a middleware that injects AclService into the request context
func WithContext(service AclGateService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(NewContext(r.Context(), service)))
		})
	}
}

// FromContext retrieves AclService from the context
func FromContext(ctx context.Context) (AclGateService, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context is nil")
	}
	service, ok := ctx.Value(aclServiceKey).(AclGateService)
	if !ok {
		return nil, fmt.Errorf("expected AclService in context, but not found or invalid type")
	}
	return service, nil
}

// MustFromContext retrieves AclService from the context and panics if not found
func MustFromContext(ctx context.Context) AclGateService {
	service, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}
	return service
}

// NewContext creates a new context with AclService
func NewContext(ctx context.Context, service AclGateService) context.Context {
	return context.WithValue(ctx, aclServiceKey, service)
}

// CheckPermission is a helper function to check permissions using the service from context
func CheckPermission(ctx context.Context, resourceType, resourceId, subjectType, subjectId, relation string) (bool, error) {
	service, err := FromContext(ctx)
	if err != nil {
		return false, err
	}
	return service.Check(ctx, resourceType, resourceId, subjectType, subjectId, relation)
}
