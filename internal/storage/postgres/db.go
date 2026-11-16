package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	// 🔥 THIS is the line that actually connects and verifies Postgres is alive!
	if err := db.Ping(); err != nil {
		db.Close() // Clean up pool if ping fails
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Successfully connected to Postgres!")
	return db, nil

}

func GetEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
