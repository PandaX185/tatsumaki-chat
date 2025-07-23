package config

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DbInstance *sqlx.DB

func InitDb() (*sqlx.DB, error) {
	if DbInstance == nil {
		connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"), os.Getenv("DB_HOST"))
		db, err := sqlx.Connect("postgres", connStr)
		if err != nil {
			return nil, err
		}
		DbInstance = db
	}
	return DbInstance, nil
}
