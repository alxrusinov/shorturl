package store

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5"
)

type DbStore struct {
	db *sql.DB
}

func (store *DbStore) GetLink(key string) (string, error) {
	var s string
	err := store.db.QueryRow("SELECT original FROM links WHERE short = $1", key).Scan(&s)

	if err != nil {
		return "", err
	}

	return s, nil
}

func (store *DbStore) SetLink(key, link string) error {
	dbQuery := `INSERT INTO links (short, original)
				VALUES($1, $2);
				`

	err := store.db.QueryRow(dbQuery, key, link)

	if err != nil {
		return err.Err()
	}

	return nil
}

func (store *DbStore) Ping() error {
	err := store.db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func CreateDbStore(dbPath string) Store {
	db, err := sql.Open("pgx", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	return &DbStore{
		db: db,
	}
}
