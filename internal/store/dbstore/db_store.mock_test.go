package dbstore

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestNewDBStoreMock(t *testing.T) {
	db, mock, _ := sqlmock.New(sqlmock.MonitorPingsOption(true))
	tests := []struct {
		name    string
		want    *DBStore
		want1   sqlmock.Sqlmock
		wantErr bool
	}{
		{
			name: "1# success",
			want: &DBStore{
				db: db,
			},
			want1:   mock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := NewDBStoreMock()
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDBStoreMock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
