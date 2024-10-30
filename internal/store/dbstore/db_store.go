package dbstore

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// DBStore implements Store interface for data base
type DBStore struct {
	db          *sql.DB
	insertQuery *sql.Stmt
}

// NewDBStore initializes and returns data base store instance
func NewDBStore(dbPath string) *DBStore {
	store := &DBStore{}

	db, err := sql.Open("pgx", dbPath)

	if err != nil {
		log.Fatal(err)
	}

	store.db = db

	err = store.createTable()

	if err != nil {
		log.Fatal(err)
	}

	insertQueryString := `INSERT INTO links (short, original, correlation_id, user_id)
				VALUES ($1, $2, $3, $4)
				RETURNING short, original, correlation_id, user_id;
				`

	insertQuery, err := db.Prepare(insertQueryString)

	if err != nil {
		log.Fatal(err)
	}

	store.insertQuery = insertQuery

	return store
}

// CloseConnection is a method closed db connection
func CloseConnection(db *sql.DB) {
	defer db.Close()

}
