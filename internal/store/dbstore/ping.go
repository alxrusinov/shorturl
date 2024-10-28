package dbstore

func (store *DBStore) Ping() error {
	err := store.db.Ping()

	if err != nil {
		return err
	}

	return nil
}
