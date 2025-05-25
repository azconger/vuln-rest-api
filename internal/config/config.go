package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Security SecurityConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
	WeakImplementation bool // Intentionally vulnerable
	PredictableTokens  bool // Intentionally vulnerable
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	RateLimitRequests int           // 0 means unlimited (intentionally vulnerable)
	RateLimitDuration time.Duration // Duration for rate limiting window
	CORSOrigins       []string      // Intentionally permissive
	DebugMode         bool          // Intentionally vulnerable
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "localhost"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 5*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "vuln_rest_api"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:          getEnv("JWT_SECRET_KEY", "weak-secret-key"), // Intentionally weak
			AccessTokenExpiry:  getDurationEnv("JWT_ACCESS_TOKEN_EXPIRY", 15*time.Minute),
			RefreshTokenExpiry: getDurationEnv("JWT_REFRESH_TOKEN_EXPIRY", 24*time.Hour),
			WeakImplementation: true, // Intentionally vulnerable
			PredictableTokens:  true, // Intentionally vulnerable
		},
		Security: SecurityConfig{
			RateLimitRequests: getIntEnv("RATE_LIMIT_REQUESTS", 0), // 0 means unlimited (intentionally vulnerable)
			RateLimitDuration: getDurationEnv("RATE_LIMIT_DURATION", 1*time.Minute),
			CORSOrigins:       getStringSliceEnv("CORS_ORIGINS", []string{"*"}), // Intentionally permissive
			DebugMode:         true,                                             // Intentionally vulnerable
		},
	}
}

// Helper functions to get environment variables with defaults
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getStringSliceEnv(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return []string{value} // Simple implementation, could be enhanced
	}
	return defaultValue
}
