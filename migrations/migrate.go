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
		createMessagesTrigger(db),
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
		// rollbackCreateUsersTable(db),
		// rollbackCreateUsersChatsTable(db),
		// rollbackCreateChatsTable(db),
		// rollbackCreateMessagesTable(db),
	}

	for _, err := range errs {
		if err != nil {
			return err
		}
	}
	return nil
}
