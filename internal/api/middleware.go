package api

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Key type for context
type contextKey string

const userCtxKey contextKey = "sub"

// Validates the JWT and injects user ID into the request context
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil || !token.Valid {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Extract sub from claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		subFloat, ok := claims["sub"].(float64)
		if !ok {
			http.Error(w, "Invalid sub in token", http.StatusUnauthorized)
			return
		}

		userID := int(subFloat)

		// Inject user ID into request context using "sub" as key
		ctx := context.WithValue(r.Context(), userCtxKey, userID)

		// Call next handler
		next(w, r.WithContext(ctx))
	}
}

// Extracts user ID from context
func userIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userCtxKey).(int)
	return userID, ok
}
