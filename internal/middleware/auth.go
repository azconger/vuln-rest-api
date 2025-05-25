package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Vulnerable JWT implementation with weak secret
var jwtSecret = []byte("very_secret_key_123")

// AuthMiddleware implements a vulnerable authentication middleware
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Vulnerable: Weak token extraction
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		// Vulnerable: No proper Bearer token validation
		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Vulnerable: Weak token validation
		// Vulnerable: No token expiration check
		// Vulnerable: No token issuer validation
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Vulnerable: No algorithm validation
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Vulnerable: No role-based access control
		next.ServeHTTP(w, r)
	}
}
