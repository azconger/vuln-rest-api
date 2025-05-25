package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

// User represents a user in the system
// @Description User information
type User struct {
	// User ID
	// required: true
	ID int `json:"id" example:"1"`
	// Username
	// required: true
	Username string `json:"username" example:"admin"`
	// Email address
	// required: true
	Email string `json:"email" example:"admin@example.com"`
	// User role
	// required: true
	Role string `json:"role" example:"admin"`
}

// HandleGetUsers implements a vulnerable GET /users endpoint
// @Summary Get users
// @Description Get list of users with optional filtering
// @Tags users
// @Accept json
// @Produce json
// @Param query query string false "SQL WHERE clause for filtering"
// @Security BearerAuth
// @Success 200 {array} User
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Server error"
// @Router /users [get]
func HandleGetUsers(w http.ResponseWriter, r *http.Request) {
	// Vulnerable: SQL injection through query parameter
	query := r.URL.Query().Get("query")

	// Vulnerable: Direct string concatenation for SQL query
	// Vulnerable: No input validation
	// Vulnerable: No parameterized queries
	sqlQuery := "SELECT id, username, email, role FROM users"
	if query != "" {
		sqlQuery += " WHERE " + query
	}

	// Get database connection details from environment variables
	dbHost := getEnvOrDefault("DB_HOST", "localhost")
	dbPort := getEnvOrDefault("DB_PORT", "5432")
	dbUser := getEnvOrDefault("DB_USER", "postgres")
	dbPassword := getEnvOrDefault("DB_PASSWORD", "postgres")
	dbName := getEnvOrDefault("DB_NAME", "vuln_db")

	// Vulnerable: Hardcoded database credentials in connection string
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		http.Error(w, "Database connection error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Vulnerable: No error handling for SQL query
	rows, err := db.Query(sqlQuery)
	if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role); err != nil {
			http.Error(w, "Error scanning results", http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	// Vulnerable: No pagination
	// Vulnerable: No rate limiting
	// Vulnerable: No access control
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
