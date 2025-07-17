package migrations

import (
	"github.com/jmoiron/sqlx"
)

func createUsersTable(db *sqlx.DB) error {
	query := `
	create table if not exists users (
		id serial primary key,
		full_name text not null,
		user_name text not null unique,
		password text not null,

		created_at timestamp default now(),
		updated_at timestamp default now(),
		deleted_at timestamp default null
	)
	`

	_, err := db.Exec(query)
	return err
}

func rollbackCreateUsersTable(db *sqlx.DB) error {
	query := `drop table if exists users`

	_, err := db.Exec(query)
	return err
}
