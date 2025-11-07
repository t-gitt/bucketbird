package middleware

import (
	"context"
	"net/http"
	"strings"

	"bucketbird/backend/internal/repository"
	"bucketbird/backend/internal/service"

	"github.com/google/uuid"
)

type contextKey string

const UserContextKey contextKey = "user"

// Auth middleware extracts and validates JWT token, adds user to context
func Auth(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract token from Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				// Also check for lowercase (some clients send lowercase)
				authHeader = r.Header.Get("authorization")
			}

			if authHeader == "" {
				http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			// Check Bearer prefix
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"Invalid authorization header"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			token := parts[1]

			// Validate token
			user, err := authService.ValidateAccessToken(r.Context(), token)
			if err != nil {
				http.Error(w, `{"error":"Invalid token"}`, http.StatusUnauthorized)
				w.Header().Set("Content-Type", "application/json")
				return
			}

			// Add user to context
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserFromContext extracts user from context
func GetUserFromContext(ctx context.Context) (*repository.User, bool) {
	user, ok := ctx.Value(UserContextKey).(*repository.User)
	return user, ok
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		return uuid.Nil, false
	}
	return user.ID, true
}

// DemoReadOnly middleware blocks write operations for demo users
func DemoReadOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := GetUserFromContext(r.Context())
		if ok && user.IsDemo {
			// Allow only GET and OPTIONS requests for demo users
			if r.Method != http.MethodGet && r.Method != http.MethodOptions {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"error":"Demo users have read-only access. Please sign up for a full account to perform this action."}`, http.StatusForbidden)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
