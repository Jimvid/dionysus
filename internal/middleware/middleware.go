package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/jimvid/dionysus/internal/jwt"
)

// Key type for storing claims in context
type contextKey string

const ClaimsContextKey contextKey = "claims"

func ValidateJWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := extractTokenFromHeader(r)
		if tokenString == "" {
			http.Error(w, "Missing or malformed Authorization header", http.StatusUnauthorized)
			return
		}

		claims, err := jwt.VerifyToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check expiration
		expires, ok := claims["expires"].(float64)
		if !ok || time.Now().Unix() > int64(expires) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Add claims to context
		ctx := context.WithValue(r.Context(), ClaimsContextKey, claims)

		// next
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func extractTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}
