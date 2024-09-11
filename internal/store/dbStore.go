package store

import (
	"context"
	"database/sql"
	"io"
	"log"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type DBStore struct {
	db          *sql.DB
	deleteQuery *sql.Stmt
	insertQuery *sql.Stmt
}

func (store *DBStore) GetLink(arg *StoreRecord) (*StoreRecord, error) {
	var s string
	err := store.db.QueryRowContext(context.Background(), "SELECT original FROM links WHERE short = $1 and is_deleted = FALSE", arg.ShortLink).Scan(&s)

	if err != nil {
		return nil, err
	}

	arg.OriginalLink = s

	return arg, nil
}

func (store *DBStore) SetLink(arg *StoreRecord) (*StoreRecord, error) {
	var err error
	dbQuery := `INSERT INTO links (short, original, correlation_id, user_id)
				VALUES ($1, $2, $3, $4);
				`

	selectQuery := `SELECT short FROM links WHERE original = $1 `

	_, err = store.db.ExecContext(context.Background(), dbQuery, arg.ShortLink, arg.OriginalLink, arg.CorrelationID, arg.UUID)

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

	stmt := tx.Stmt(store.insertQuery)

	defer stmt.Close()

	response := make([]*StoreRecord, 0)

	for _, val := range arg {
		res := &StoreRecord{}
		err := stmt.QueryRowContext(context.Background(), val.ShortLink, val.OriginalLink, val.CorrelationID, val.UUID).Scan(&res.ShortLink, &res.OriginalLink, &res.CorrelationID, &res.UUID)

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

func (store *DBStore) GetLinks(userID string) ([]StoreRecord, error) {
	rows, err := store.db.QueryContext(context.Background(), "SELECT user_id, short, original, correlation_id, is_deleted FROM links WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []StoreRecord

	for rows.Next() {
		var row StoreRecord

		if err := rows.Scan(&row.UUID, &row.ShortLink, &row.OriginalLink, &row.CorrelationID, &row.Deleted); err != nil {
			return nil, err
		}

		result = append(result, row)

	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil

}

func (store *DBStore) DeleteLinks(shorts [][]StoreRecord) error {
	tx, err := store.db.Begin()

	if err != nil {
		return err
	}

	defer tx.Rollback()

	preparedShorts := make([]string, len(shorts))
	preparedIDs := make([]string, len(shorts))

	for _, val := range shorts {
		userID := val[0].UUID
		preparedIDs = append(preparedIDs, userID)

		for _, shortLink := range val {
			preparedShorts = append(preparedShorts, shortLink.ShortLink)
		}
	}

	shortsPlaceholders := strings.Join(preparedShorts, ", ")
	userIDsPlaceholders := strings.Join(preparedIDs, ", ")

	stmt := tx.Stmt(store.deleteQuery)

	defer stmt.Close()

	_, err = stmt.ExecContext(context.Background(), userIDsPlaceholders, shortsPlaceholders)

	if err != nil {
		return err
	}

	return nil

}

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

func CreateDBStore(dbPath string) Store {
	db, err := sql.Open("pgx", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	insertQueryString := `INSERT INTO links (short, original, correlation_id, user_id)
				VALUES ($1, $2, $3, $4)
				RETURNING short, original, correlation_id, user_id;
				`

	deleteQueryString := `UPDATE links SET is_deleted = TRUE WHERE user_id in ($1) and short IN ($2)`

	insertQuery, err := db.Prepare(insertQueryString)

	if err != nil {
		log.Fatal(err)
	}

	deleteQuery, err := db.Prepare(deleteQueryString)

	if err != nil {
		log.Fatal(err)
	}

	store := &DBStore{
		db:          db,
		insertQuery: insertQuery,
		deleteQuery: deleteQuery,
	}

	err = store.createTable()

	if err != nil {
		log.Fatal(err)
	}

	return store
}

func CloseConnection(db *sql.DB) {
	defer db.Close()

}
