package database

import (
	"database/sql"
	"embed"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql seed.sql
var sqlFS embed.FS

var (
	db *sql.DB
)

// InitDB initializes the database connection
func InitDB() error {
	var err error

	// Use in-memory database by default
	dbPath := "file::memory:?cache=shared"

	// If DB_PATH is set, use that instead
	if path := os.Getenv("DB_PATH"); path != "" {
		// Ensure directory exists
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		dbPath = path
	}

	// Open database connection
	db, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}

	// Read and execute schema
	schema, err := sqlFS.ReadFile("schema.sql")
	if err != nil {
		return err
	}

	// Execute schema
	if _, err := db.Exec(string(schema)); err != nil {
		return err
	}

	// Read and execute seed data
	seed, err := sqlFS.ReadFile("seed.sql")
	if err != nil {
		return err
	}

	// Execute seed data
	if _, err := db.Exec(string(seed)); err != nil {
		return err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return err
	}

	log.Printf("Database initialized successfully at %s", dbPath)
	return nil
}

// GetDB returns the database connection
func GetDB() *sql.DB {
	return db
}

// CloseDB closes the database connection
func CloseDB() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
