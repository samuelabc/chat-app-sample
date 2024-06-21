package auth

import (
	"context"
	"net/http"
	"strings"
)

// JWTMiddleware is a middleware function that validates JWT tokens
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userId, username, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userId", userId)
		ctx = context.WithValue(ctx, "username", username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
