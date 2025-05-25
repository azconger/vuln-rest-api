package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/azconger/vuln-rest-api/internal/handlers"
)

func TestGetUsersEndpoint(t *testing.T) {
	// Test case 1: Basic request without query
	req := httptest.NewRequest("GET", "/api/v1/users", nil)
	w := httptest.NewRecorder()

	handlers.HandleGetUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Test case 2: Request with SQL injection attempt
	req = httptest.NewRequest("GET", "/api/v1/users?query=1%3D1%3BDROP%20TABLE%20users%3B--", nil)
	w = httptest.NewRecorder()

	handlers.HandleGetUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Test case 3: Request with malicious query
	req = httptest.NewRequest("GET", "/api/v1/users?query=username%3D%27admin%27%20OR%20%271%27%3D%271%27", nil)
	w = httptest.NewRecorder()

	handlers.HandleGetUsers(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
