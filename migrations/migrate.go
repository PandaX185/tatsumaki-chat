package migrations

import (
	"github.com/jmoiron/sqlx"
)

func Up(db *sqlx.DB) error {
	var err error
	if err = createUsersTable(db); err != nil {
		return err
	}

	return nil
}

func Down(db *sqlx.DB) error {
	var err error
	if err = rollbackCreateUsersTable(db); err != nil {
		return err
	}

	return nil
}
