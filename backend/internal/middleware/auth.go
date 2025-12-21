package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/homelab/filemanager/internal/service"
)

// ContextKey is a type for context keys to avoid collisions
type ContextKey string

const (
	// UserClaimsKey is the context key for user claims
	UserClaimsKey ContextKey = "userClaims"
)

// JWTAuth creates a middleware that validates JWT tokens
func JWTAuth(authService service.AuthService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeAuthError(w, "Missing authorization header", http.StatusUnauthorized)
				return
			}

			// Check Bearer prefix
			if !strings.HasPrefix(authHeader, "Bearer ") {
				writeAuthError(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// Extract token
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == "" {
				writeAuthError(w, "Missing token", http.StatusUnauthorized)
				return
			}

			// Validate token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				switch err {
				case service.ErrTokenExpired:
					writeAuthError(w, "Token expired", http.StatusUnauthorized)
				case service.ErrInvalidToken:
					writeAuthError(w, "Invalid token", http.StatusUnauthorized)
				default:
					writeAuthError(w, "Authentication failed", http.StatusUnauthorized)
				}
				return
			}

			// Add claims to context
			ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserClaims retrieves user claims from the request context
func GetUserClaims(ctx context.Context) (*service.Claims, bool) {
	claims, ok := ctx.Value(UserClaimsKey).(*service.Claims)
	return claims, ok
}

// writeAuthError writes an authentication error response
func writeAuthError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(`{"error":"` + message + `","code":"UNAUTHORIZED"}`))
}
