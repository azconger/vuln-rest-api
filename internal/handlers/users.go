package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	_ "github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// HandleGetUsers implements a vulnerable GET /users endpoint
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

	// Vulnerable: Hardcoded database credentials
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/vuln_db?sslmode=disable")
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
