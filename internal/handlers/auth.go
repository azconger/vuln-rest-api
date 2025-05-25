package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Vulnerable JWT implementation with weak secret and predictable token generation
var jwtSecret = []byte("very_secret_key_123") // Weak secret key

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// Vulnerable: Hardcoded credentials
var validCredentials = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

// HandleLogin implements a vulnerable login endpoint
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Vulnerable: Weak password comparison
	// Vulnerable: No password hashing
	// Vulnerable: No rate limiting
	// Vulnerable: No account lockout
	expectedPassword, exists := validCredentials[req.Username]
	if !exists || expectedPassword != req.Password {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Vulnerable: Predictable token generation
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	})

	// Vulnerable: Weak signing method (HS256)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Vulnerable: Predictable refresh token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": req.Username,
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
	})

	refreshTokenString, err := refreshToken.SignedString(jwtSecret)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	response := TokenResponse{
		AccessToken:  tokenString,
		TokenType:    "Bearer",
		ExpiresIn:    86400, // 24 hours in seconds
		RefreshToken: refreshTokenString,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleToken implements a vulnerable OAuth ROPC token endpoint
func HandleToken(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: Accepts any client credentials
	// Vulnerable: No rate limiting
	HandleLogin(w, r)
}

// HandleRefresh implements a vulnerable token refresh endpoint
func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: No validation of refresh token
	// Vulnerable: Reuses the same token generation logic
	HandleLogin(w, r)
}

// HandleLogout implements a vulnerable logout endpoint
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: No token invalidation
	// Vulnerable: No session management
	w.WriteHeader(http.StatusOK)
}
