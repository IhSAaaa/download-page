package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

// InitDatabase initializes the database connection
func InitDatabase() error {
	// Get database connection details from environment
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "url_shortener")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLMode)

	// Open database connection
	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	// Test the connection
	if err := DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)

	log.Println("Database connection established successfully")
	return nil
}

// CloseDatabase closes the database connection
func CloseDatabase() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// CreateTables creates all necessary tables
func CreateTables() error {
	// Create URLs table
	urlsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id UUID PRIMARY KEY,
		original_url TEXT NOT NULL,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		custom_code VARCHAR(50) UNIQUE,
		title VARCHAR(255),
		description TEXT,
		user_id UUID,
		is_active BOOLEAN DEFAULT true,
		expires_at TIMESTAMP,
		click_count BIGINT DEFAULT 0,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create clicks table
	clicksTable := `
	CREATE TABLE IF NOT EXISTS clicks (
		id UUID PRIMARY KEY,
		url_id UUID NOT NULL REFERENCES urls(id) ON DELETE CASCADE,
		ip_address INET NOT NULL,
		user_agent TEXT,
		referer TEXT,
		country VARCHAR(100),
		city VARCHAR(100),
		device VARCHAR(50),
		browser VARCHAR(50),
		os VARCHAR(50),
		clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	// Create indexes
	indexes := []string{
		"CREATE INDEX IF NOT EXISTS idx_urls_short_code ON urls(short_code);",
		"CREATE INDEX IF NOT EXISTS idx_urls_custom_code ON urls(custom_code);",
		"CREATE INDEX IF NOT EXISTS idx_urls_user_id ON urls(user_id);",
		"CREATE INDEX IF NOT EXISTS idx_clicks_url_id ON clicks(url_id);",
		"CREATE INDEX IF NOT EXISTS idx_clicks_clicked_at ON clicks(clicked_at);",
		"CREATE INDEX IF NOT EXISTS idx_urls_created_at ON urls(created_at);",
	}

	// Execute table creation
	if _, err := DB.Exec(urlsTable); err != nil {
		return fmt.Errorf("failed to create urls table: %v", err)
	}

	if _, err := DB.Exec(clicksTable); err != nil {
		return fmt.Errorf("failed to create clicks table: %v", err)
	}

	// Execute indexes
	for _, index := range indexes {
		if _, err := DB.Exec(index); err != nil {
			return fmt.Errorf("failed to create index: %v", err)
		}
	}

	log.Println("Database tables created successfully")
	return nil
}

// getEnv gets environment variable with fallback
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// GetDB returns the database instance
func GetDB() *sql.DB {
	return DB
} 