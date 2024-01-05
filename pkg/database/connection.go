package database

import (
	"database/sql"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() (*sql.DB, error) {
	connectionStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return nil, err
	}

	return db, nil
}
