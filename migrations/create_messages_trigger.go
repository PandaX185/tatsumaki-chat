package migrations

import "github.com/jmoiron/sqlx"

func createMessagesTrigger(tx *sqlx.DB) error {
	function := `
	create function notify_new_message() returns trigger as $$
	begin
		perform pg_notify('message_sent', row_to_json(new)::text);
		return new;
	end;
	$$ language plpgsql;
	`

	tx.MustExec(function)

	trigger := `
	create or replace trigger on_message_sent 
	after insert on public.messages 
	for each row execute function notify_new_message();
	`

	tx.MustExec(trigger)

	return nil
}
