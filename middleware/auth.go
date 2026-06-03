package middleware

import (
	"net/http"
	"strings"
)

const (
	authHeaderKey = "Authorization"
	authPrefix    = "Bearer "
)

// IsAuthed middleware checks for a valid Bearer token.
// Returns 401 Unauthorized if missing/invalid.
func IsAuthed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Check if Authorization header exists.
		authHeader := r.Header.Get(authHeaderKey)
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "missing Authorization header"}`))
			return
		}

		// 2. Extract token (remove "Bearer " prefix).
		token := strings.TrimPrefix(authHeader, authPrefix)
		if token == authHeader { // Prefix wasn't found.
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "invalid Authorization format: use 'Bearer <token>'"}`))
			return
		}

		// 3. Validate token (example: check JWT using a library like github.com/golang-jwt/jwt).
		if !isTokenValid(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "invalid token"}`))
			return
		}

		// 4. Token is valid; proceed.
		next.ServeHTTP(w, r)
	})
}

// Placeholder: Replace with actual token validation logic (e.g., JWT parsing).
func isTokenValid(token string) bool {
	// Example: Parse and verify a JWT token.
	// In production, use a library like:
	// token, err := jwt.Parse(token, func(*jwt.Token) (interface{}, error) {
	//     return []byte(secretKey), nil
	// })
	// if err != nil || !token.Valid { return false }
	return token != "" // Dummy validation for example.
}
