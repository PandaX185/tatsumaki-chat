package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() {
	connectionStr := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
}
