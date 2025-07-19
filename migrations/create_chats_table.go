package migrations

import "github.com/jmoiron/sqlx"

func createChatsTable(db *sqlx.DB) error {
	query := `
	create table if not exists chats (
		id serial primary key,
		chat_name text not null,
		chat_owner int not null,
		created_at timestamp not null default now(),

		constraint fk_chat_owner foreign key (chat_owner) references users(id) on delete cascade
	)
	`

	_, err := db.Exec(query)
	return err
}

func rollbackCreateChatsTable(db *sqlx.DB) error {
	query := `drop table if exists chats`

	_, err := db.Exec(query)
	return err
}
