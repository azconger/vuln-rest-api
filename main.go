package main

import (
	"log"
	"net/http"
	"os"

	"github.com/azconger/vuln-rest-api/internal/handlers"
	"github.com/azconger/vuln-rest-api/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize router
	r := mux.NewRouter()

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/login", handlers.HandleLogin).Methods("POST")
	api.HandleFunc("/auth/token", handlers.HandleToken).Methods("POST")
	api.HandleFunc("/auth/refresh", handlers.HandleRefresh).Methods("POST")
	api.HandleFunc("/auth/logout", handlers.HandleLogout).Methods("POST")

	// User routes (protected)
	api.HandleFunc("/users", middleware.AuthMiddleware(handlers.HandleGetUsers)).Methods("GET")

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
