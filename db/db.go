package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (
	// DB is the global database connection pool
	DB *sql.DB
)

// User represents a user in the system
// type User struct {
// 	ID           int    `json:"id"`
// 	Email        string `json:"email"`
// 	PasswordHash string `json:"-"`
// 	FullName     string `json:"full_name"`
// }

// InitDB initializes the database connection
func InitDB() error {
	if DB != nil {
		return nil // Already initialized
	}

	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}

	// Test the connection
	if err = DB.Ping(); err != nil {
		DB.Close() // Clean up on failure
		return fmt.Errorf("error connecting to database: %w", err)
	}

	// Initialize schema
	if err = createTables(); err != nil {
		DB.Close() // Clean up on failure
		return fmt.Errorf("error creating tables: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

// createTables ensures all required tables exist
func createTables() error {
	query := `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            email VARCHAR(255) UNIQUE NOT NULL,
            password_hash VARCHAR(255) NOT NULL,
            full_name VARCHAR(255) NOT NULL
        );`

	_, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("error creating users table: %w", err)
	}
	return nil
}
