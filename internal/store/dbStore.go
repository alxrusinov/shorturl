package store

import (
	"database/sql"
	"io"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type DBStore struct {
	db *sql.DB
}

func (store *DBStore) GetLink(arg *StoreArgs) (string, error) {
	var s string
	err := store.db.QueryRow("SELECT original FROM links WHERE short = $1", arg.ShortLink).Scan(&s)

	if err != nil {
		return "", err
	}

	return s, nil
}

func (store *DBStore) SetLink(arg *StoreArgs) error {
	dbQuery := `INSERT INTO links (short, original, correlation_id)
				VALUES ($1, $2, $3);
				`

	_, err := store.db.Exec(dbQuery, arg.ShortLink, arg.OriginalLink, arg.CorrelationID)

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

func (store *DBStore) SetBatchLink(arg []*StoreArgs) ([]*StoreArgs, error) {
	tx, err := store.db.Begin()

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	insertQuery := `INSERT INTO links (short, original, correlation_id)
				VALUES ($1, $2, $3)
				RETURNING short, original, correlation_id;
				`

	stmt, err := tx.Prepare(insertQuery)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	response := make([]*StoreArgs, len(arg))

	for _, val := range arg {
		res := &StoreArgs{}
		err := stmt.QueryRow(val.ShortLink, val.OriginalLink, val.CorrelationID).Scan(&res.ShortLink, &res.OriginalLink, &res.CorrelationID)

		if err != nil && err != io.EOF {
			return nil, err
		}

		response = append(response, res)

	}

	err = tx.Commit()

	if err != nil {
		return nil, err
	}

	return response, nil
}

func CreateDBStore(dbPath string) Store {
	db, err := sql.Open("pgx", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	initialQuery := `CREATE TABLE IF NOT EXISTS links (
		id SERIAL PRIMARY KEY,
		short TEXT,
		original TEXT,
		correlation_id TEXT
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
