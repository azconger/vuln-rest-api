package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Vulnerable JWT implementation with weak secret and predictable token generation
var jwtSecret = []byte("very_secret_key_123") // Weak secret key

// LoginRequest represents the login request body
// @Description Login request payload
type LoginRequest struct {
	// Username for login
	// required: true
	Username string `json:"username" example:"admin"`
	// Password for login
	// required: true
	Password string `json:"password" example:"admin123"`
}

// TokenResponse represents the token response
// @Description Token response payload
type TokenResponse struct {
	// JWT access token
	// required: true
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	// Token type
	// required: true
	TokenType string `json:"token_type" example:"Bearer"`
	// Token expiration in seconds
	// required: true
	ExpiresIn int64 `json:"expires_in" example:"86400"`
	// JWT refresh token
	// required: true
	RefreshToken string `json:"refresh_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// Vulnerable: Hardcoded credentials
var validCredentials = map[string]string{
	"admin": "admin123",
	"user":  "user123",
}

// HandleLogin implements a vulnerable login endpoint
// @Summary Login to get JWT token
// @Description Login with username and password to get JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Server error"
// @Router /auth/login [post]
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
// @Summary Get OAuth token
// @Description Get OAuth token using Resource Owner Password Credentials flow
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Server error"
// @Router /auth/token [post]
func HandleToken(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: Accepts any client credentials
	// Vulnerable: No rate limiting
	HandleLogin(w, r)
}

// HandleRefresh implements a vulnerable token refresh endpoint
// @Summary Refresh JWT token
// @Description Refresh JWT token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login credentials"
// @Success 200 {object} TokenResponse
// @Failure 400 {string} string "Invalid request"
// @Failure 401 {string} string "Invalid credentials"
// @Failure 500 {string} string "Server error"
// @Router /auth/refresh [post]
func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: No validation of refresh token
	// Vulnerable: Reuses the same token generation logic
	HandleLogin(w, r)
}

// HandleLogout implements a vulnerable logout endpoint
// @Summary Logout
// @Description Logout and invalidate token
// @Tags auth
// @Success 200 {string} string "OK"
// @Router /auth/logout [post]
func HandleLogout(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: No token invalidation
	// Vulnerable: No session management
	w.WriteHeader(http.StatusOK)
}
