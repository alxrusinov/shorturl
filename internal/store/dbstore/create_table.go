package dbstore

import "context"

func (store *DBStore) createTable() error {
	initialQuery := `CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		user_id TEXT,
		short TEXT,
		original TEXT UNIQUE,
		correlation_id TEXT,
		is_deleted BOOLEAN NOT NULL DEFAULT FALSE
	);`

	_, err := store.db.ExecContext(context.Background(), initialQuery)

	if err != nil {
		return err
	}

	return nil
}
