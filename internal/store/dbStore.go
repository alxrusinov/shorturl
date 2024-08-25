package store

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type DBStore struct {
	db *sql.DB
}

func (store *DBStore) GetLink(arg *StoreRecord) (*StoreRecord, error) {
	var s string
	err := store.db.QueryRowContext(context.Background(), "SELECT original FROM links WHERE short = $1", arg.ShortLink).Scan(&s)

	if err != nil {
		return nil, err
	}

	arg.OriginalLink = s

	return arg, nil
}

func (store *DBStore) SetLink(arg *StoreRecord) (*StoreRecord, error) {
	var err error
	dbQuery := `INSERT INTO links (short, original, correlation_id)
				VALUES ($1, $2, $3);
				`

	selectQuery := `SELECT short FROM links WHERE original = $1 `

	_, err = store.db.ExecContext(context.Background(), dbQuery, arg.ShortLink, arg.OriginalLink, arg.CorrelationID)

	if err != nil {
		if dbErr, ok := err.(*pgconn.PgError); ok {
			if dbErr.Code == pgerrcode.UniqueViolation {

				err := store.db.QueryRowContext(context.Background(), selectQuery, arg.OriginalLink).Scan(&arg.ShortLink)

				if err != nil {
					return nil, err
				}

				return arg, &DuplicateValueError{Err: dbErr}
			}
		}
		return nil, err
	}

	return arg, nil
}

func (store *DBStore) Ping() error {
	err := store.db.Ping()

	if err != nil {
		return err
	}

	return nil
}

func (store *DBStore) SetBatchLink(arg []*StoreRecord) ([]*StoreRecord, error) {
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

	response := make([]*StoreRecord, 0)

	for _, val := range arg {
		res := &StoreRecord{}
		err := stmt.QueryRowContext(context.Background(), val.ShortLink, val.OriginalLink, val.CorrelationID).Scan(&res.ShortLink, &res.OriginalLink, &res.CorrelationID)

		if err != nil && !errors.Is(err, io.EOF) {
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
		original TEXT UNIQUE,
		correlation_id TEXT
	);`

	_, err = db.ExecContext(context.Background(), initialQuery)

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
