package storage

import (
	"database/sql"
	"embed"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"os"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func InitDB() (*sql.DB, error) {
	// Load environment variables
	godotenv.Load() // intentionally ignoring error

	// Check for required environment variables (excluding password)
	requiredEnvVars := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_NAME"}
	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			return nil, fmt.Errorf("required environment variable %s is not set", envVar)
		}
	}

	// Build connection string based on whether password exists
	var connStr string
	if os.Getenv("DB_PASSWORD") == "" {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
		)
	} else {
		connStr = fmt.Sprintf(
			"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_PASSWORD"),
			os.Getenv("DB_NAME"),
		)
	}

	// Open a connection to the database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the database: %v", err)
	}

	// Testing the connection
	if err := db.Ping(); err != nil {
		db.Close() // Clean up the connection if ping fails
		return nil, fmt.Errorf("error pinging the database: %v", err)
	}

	// Set up goose with our embedded migrations
	goose.SetBaseFS(embedMigrations)

	// Run database migrations
	if err := goose.Up(db, "migrations"); err != nil {
		db.Close() // Clean up the connection if migrations fail
		return nil, fmt.Errorf("error running migrations: %v", err)
	}

	return db, nil
}

