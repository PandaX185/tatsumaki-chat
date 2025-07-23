package migrations

import (
	"github.com/jmoiron/sqlx"
)

func createMessagesTable(db *sqlx.DB) error {
	query := `
	create table if not exists messages (
		id serial primary key,
		cid integer not null,
		sender_id integer not null,
		content text not null default '',

		created_at timestamp not null default now(),
		updated_at timestamp not null default now(),
		constraint fk_cid foreign key (cid) references chats (id) on delete cascade,
		constraint fk_sender foreign key (sender_id) references users (id) on delete cascade
	)
	`

	_, err := db.Exec(query)
	return err
}

func rollbackCreateMessagesTable(db *sqlx.DB) error {
	query := `drop table if exists messages`

	_, err := db.Exec(query)
	return err
}
