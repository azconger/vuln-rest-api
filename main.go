package main

import (
	"log"
	"net/http"
	"os"

	"github.com/azconger/vuln-rest-api/docs" // Import generated Swagger docs
	"github.com/azconger/vuln-rest-api/internal/database"
	"github.com/azconger/vuln-rest-api/internal/handlers"
	"github.com/azconger/vuln-rest-api/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Vulnerable REST API
// @version 1.0
// @description A deliberately vulnerable REST API for testing and demonstration purposes.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

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

	// Command injection routes (protected)
	api.HandleFunc("/commands", middleware.AuthMiddleware(handlers.HandleCommand)).Methods("POST")
	api.HandleFunc("/commands/shell", middleware.AuthMiddleware(handlers.HandleShellCommand)).Methods("POST")
	api.HandleFunc("/commands/filter", middleware.AuthMiddleware(handlers.HandleCommandWithFilter)).Methods("POST")

	// Path traversal routes (protected)
	api.HandleFunc("/files/read", middleware.AuthMiddleware(handlers.HandleFileRead)).Methods("POST")
	api.HandleFunc("/files/download", middleware.AuthMiddleware(handlers.HandleFileDownload)).Methods("POST")
	api.HandleFunc("/files/write", middleware.AuthMiddleware(handlers.HandleFileWrite)).Methods("POST")
	api.HandleFunc("/files/list", middleware.AuthMiddleware(handlers.HandleFileList)).Methods("POST")

	// Swagger documentation
	docs.SwaggerInfo.Title = "Vulnerable REST API"
	docs.SwaggerInfo.Description = "A deliberately vulnerable REST API for testing and demonstration purposes."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	// Configure Swagger UI with authorization
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
