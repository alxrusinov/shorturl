package dbstore

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
)

func TestDBStore_DeleteLinks(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()

	type args struct {
		shorts [][]model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{shorts: [][]model.StoreRecord{{{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}}}},
			wantErr: false,
		},
		{
			name:  "2# error transaction",
			store: teststore,
			args: args{shorts: [][]model.StoreRecord{{{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}}}},
			wantErr: true,
		},
		{
			name:  "2# error",
			store: teststore,
			args: args{shorts: [][]model.StoreRecord{{{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}}}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE links").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("1"))
				mock.ExpectCommit()
			}

			if tt.name == tests[1].name {
				mock.ExpectBegin().WillReturnError(errors.New("tr err"))
			}
			if tt.name == tests[2].name {
				mock.ExpectBegin()
				mock.ExpectQuery("UPDATE links").WillReturnRows(sqlmock.NewRows([]string{"id"})).WillReturnError(errors.New("err"))
				mock.ExpectCommit()
			}

			if err := tt.store.DeleteLinks(tt.args.shorts); (err != nil) != tt.wantErr {
				t.Errorf("DBStore.DeleteLinks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
