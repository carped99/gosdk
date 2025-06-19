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
func WithContext(service AclService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(NewContext(r.Context(), service)))
		})
	}
}

// FromContext retrieves AclService from the context
func FromContext(ctx context.Context) (AclService, error) {
	service, ok := ctx.Value(aclServiceKey).(AclService)
	if !ok {
		return nil, fmt.Errorf("expected AclService in context, but not found or invalid type")
	}
	return service, nil
}

// MustFromContext retrieves AclService from the context and panics if not found
func MustFromContext(ctx context.Context) AclService {
	service, err := FromContext(ctx)
	if err != nil {
		panic(err)
	}
	return service
}

// NewContext creates a new context with AclService
func NewContext(ctx context.Context, service AclService) context.Context {
	return context.WithValue(ctx, aclServiceKey, service)
}

//// CheckPermission is a helper function to check permissions using the service from context
//func CheckPermission(ctx context.Context, userID string, resource string, action string) error {
//	service, err := FromContext(ctx)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

// GetPermissions is a helper function to get user permissions using the service from context
//func GetPermissions(ctx context.Context, userID string) ([]string, error) {
//	service, err := FromContext(ctx)
//	if err != nil {
//		return nil, err
//	}
//	return service.GetUserPermissions(ctx, userID)
//}

// RequirePermission is a middleware that checks if the user has the required permission
//func RequirePermission(resource string, action string) func(http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			// Get user ID from request context (assuming it's set by auth middleware)
//			userID := r.Context().Value("user_id")
//			if userID == nil {
//				http.Error(w, "Unauthorized", http.StatusUnauthorized)
//				return
//			}
//
//			// Check permission
//			err := CheckPermission(r.Context(), userID.(string), resource, action)
//			if err != nil {
//				log.Printf("Permission check failed: user_id=%v, resource=%v, action=%v, error=%v",
//					userID, resource, action, err)
//				http.Error(w, "Forbidden", http.StatusForbidden)
//				return
//			}
//
//			next.ServeHTTP(w, r)
//		})
//	}
//}
