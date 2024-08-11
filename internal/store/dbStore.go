package store

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStore struct {
	db *sql.DB
}

func (store *DBStore) GetLink(key string) (string, error) {
	var s string
	err := store.db.QueryRow("SELECT original FROM links WHERE short = $1", key).Scan(&s)

	if err != nil {
		return "", err
	}

	return s, nil
}

func (store *DBStore) SetLink(key, link string) error {
	dbQuery := `INSERT INTO links VALUES
				($1, $2);
				`

	_, err := store.db.Exec(dbQuery, key, link)

	if err != nil {
		return err
	}

	return nil
}

func (store *DBStore) Ping() error {
	err := store.db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func CreateDBStore(dbPath string) Store {
	db, err := sql.Open("pgx", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	initialQuery := `CREATE TABLE IF NOT EXISTS links (
		short TEXT,
		original TEXT
	);`

	_, err = db.Exec(initialQuery)

	if err != nil {
		log.Fatal(err)
	}

	return &DBStore{
		db: db,
	}
}

func CloseConnection(db *sql.DB) {
	defer db.Close()

}
