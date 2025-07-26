package migrations

import (
	"github.com/jmoiron/sqlx"
)

func Up(db *sqlx.DB) error {
	errs := []error{
		createUsersTable(db),
		createChatsTable(db),
		createUsersChatsTable(db),
		createMessagesTable(db),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}

func Down(db *sqlx.DB) error {
	errs := []error{
		// rollbackCreateMessagesTable(db),
		// rollbackCreateUsersChatsTable(db),
		// rollbackCreateChatsTable(db),
		// rollbackCreateUsersTable(db),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
