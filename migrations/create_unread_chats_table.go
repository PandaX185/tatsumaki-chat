package migrations

import "github.com/jmoiron/sqlx"

func createUnreadChatsTable(db *sqlx.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS unread_chats(
		id SERIAL PRIMARY KEY,
		uid INT NOT NULL REFERENCES users(id),
		cid INT NOT NULL REFERENCES chats(id),
		unread_count INT NOT NULL DEFAULT 0
	)
	`
	_, err := db.Exec(query)
	return err
}

func rollbackCreateUnreadChatsTable(db *sqlx.DB) error {
	query := `
	DROP TABLE IF EXISTS unread_chats
	`
	_, err := db.Exec(query)
	return err
}
