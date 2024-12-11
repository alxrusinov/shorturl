package dbstore

import (
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewDBStoreMock creates db mock instance
func NewDBStoreMock() (*DBStore, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))

	if err != nil {
		return nil, nil, err
	}

	storeMock := &DBStore{
		db: db,
	}

	return storeMock, mock, nil

}
