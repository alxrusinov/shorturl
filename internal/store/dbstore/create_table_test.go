package dbstore

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDBStore_createTableSuccess(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	tests := []struct {
		name    string
		store   *DBStore
		wantErr bool
	}{
		{
			name:    "1#success",
			store:   teststore,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("CREATE TABLE IF NOT EXISTS links").WillReturnResult(sqlmock.NewResult(1, 1))
			if err := tt.store.createTable(); (err != nil) != tt.wantErr {
				t.Errorf("DBStore.createTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDBStore_createTableError(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	tests := []struct {
		name    string
		store   *DBStore
		wantErr bool
	}{
		{
			name:    "1#success",
			store:   teststore,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("CREATE TABLE IF NOT EXISTS links").WillReturnError(errors.New("err"))
			if err := tt.store.createTable(); (err != nil) != tt.wantErr {
				t.Errorf("DBStore.createTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
