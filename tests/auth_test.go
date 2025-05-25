package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azconger/vuln-rest-api/internal/handlers"
	"github.com/azconger/vuln-rest-api/internal/middleware"
)

func TestLoginEndpoint(t *testing.T) {
	// Test case 1: Valid login request
	loginReq := handlers.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handlers.HandleLogin(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response handlers.TokenResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}

	if response.RefreshToken == "" {
		t.Error("Expected refresh token, got empty string")
	}

	// Test case 2: Invalid credentials
	loginReq = handlers.LoginRequest{
		Username: "admin",
		Password: "wrongpassword",
	}
	reqBody, _ = json.Marshal(loginReq)
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	w = httptest.NewRecorder()

	handlers.HandleLogin(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Test case 3: Invalid request body
	req = httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer([]byte("invalid json")))
	w = httptest.NewRecorder()

	handlers.HandleLogin(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestTokenEndpoint(t *testing.T) {
	// Test case 1: Valid token request
	loginReq := handlers.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/token", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handlers.HandleToken(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response handlers.TokenResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}
}

func TestRefreshEndpoint(t *testing.T) {
	// Test case 1: Valid refresh request
	loginReq := handlers.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	reqBody, _ := json.Marshal(loginReq)
	req := httptest.NewRequest("POST", "/api/v1/auth/refresh", bytes.NewBuffer(reqBody))
	w := httptest.NewRecorder()

	handlers.HandleRefresh(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response handlers.TokenResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("Expected access token, got empty string")
	}
}

func TestLogoutEndpoint(t *testing.T) {
	// Test case 1: Logout request
	req := httptest.NewRequest("POST", "/api/v1/auth/logout", nil)
	w := httptest.NewRecorder()

	handlers.HandleLogout(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProtectedRoute(t *testing.T) {
	// Test case 1: Access protected route without token
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	middleware.AuthMiddleware(handlers.HandleGetUsers)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Test case 2: Access protected route with invalid token
	req = httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w = httptest.NewRecorder()

	middleware.AuthMiddleware(handlers.HandleGetUsers)(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}

	// Test case 3: Access protected route with valid token
	// First, get a valid token
	loginReq := handlers.LoginRequest{
		Username: "admin",
		Password: "admin123",
	}
	reqBody, _ := json.Marshal(loginReq)
	loginRequest := httptest.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	loginW := httptest.NewRecorder()

	handlers.HandleLogin(loginW, loginRequest)

	var tokenResponse handlers.TokenResponse
	json.NewDecoder(loginW.Body).Decode(&tokenResponse)

	// Now use the token to access the protected route
	req = httptest.NewRequest("GET", "/api/v1/users", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)
	w = httptest.NewRecorder()

	middleware.AuthMiddleware(handlers.HandleGetUsers)(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
