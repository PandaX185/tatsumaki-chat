package migrations

import "github.com/jmoiron/sqlx"

func createUsersChatsTable(db *sqlx.DB) error {
	query := `
	create table if not exists users_chats (
		id serial primary key,
		cid integer not null,
		uid integer not null,
		created_at timestamp not null default now(),

		constraint fk_chat_id foreign key (cid) references chats(id) on delete cascade,
		constraint fk_user_id foreign key (uid) references users(id) on delete cascade,
		constraint uq_user_chat unique (cid, uid)
	)
	`

	_, err := db.Exec(query)
	return err
}

func rollbackCreateUsersChatsTable(db *sqlx.DB) error {
	query := `drop table if exists users_chats`

	_, err := db.Exec(query)
	return err
}
