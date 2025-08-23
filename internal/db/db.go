package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	// Get environment variables
	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "postgres" // デフォルト値
	}

	dbPass := os.Getenv("DB_PASSWORD")
	if dbPass == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "elmo_db" // デフォルト値
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		return nil, fmt.Errorf("DB_HOST environment variable is required")
	}

	// Cloud SQL Unix Domain Socket接続文字列を構築
	// 例: /cloudsql/project:region:instance
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s", 
		dbUser, dbPass, dbName, dbHost)

	// データベースに接続
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// 接続をテスト
	err = db.Ping()
	if err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to database %s as %s", dbName, dbUser)
	return db, nil
}