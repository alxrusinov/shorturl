package dbstore

// Ping pings data base
func (store *DBStore) Ping() error {
	err := store.db.Ping()

	if err != nil {
		return err
	}

	return nil
}
