package aclgate

import "net/http"

// WithContext creates a middleware that injects ClientService into the request context
func WithContext(service ClientService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(NewContext(r.Context(), service)))
		})
	}
}
