package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DbInstance *sql.DB

func InitDb() (*sql.DB, error) {
	if DbInstance == nil {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=require", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, err
		}
		DbInstance = db
	}
	return DbInstance, nil
}
