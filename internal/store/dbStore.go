package store

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx"
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
	dbQuery := `INSERT INTO links (short, original)
				VALUES($1, $2);
				`

	err := store.db.QueryRow(dbQuery, key, link)

	if err != nil {
		return err.Err()
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

	return &DBStore{
		db: db,
	}
}
