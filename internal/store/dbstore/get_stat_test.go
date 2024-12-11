package dbstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
)

func TestDBStore_GetStat(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	tests := []struct {
		name    string
		store   *DBStore
		want    *model.StatResponse
		wantErr bool
	}{
		{
			name:    "1# success",
			store:   teststore,
			want:    &model.StatResponse{URLS: 3, Users: 2},
			wantErr: false,
		},
		{
			name:    "2# error",
			store:   teststore,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectQuery("SELECT count\\(\\*\\) AS users, count\\(DISTINCT user_id\\) AS links FROM links").WithoutArgs().WillReturnRows(sqlmock.NewRows([]string{"urls", "users"}).AddRow("3", "2"))
			} else {
				mock.ExpectQuery("SELECT count\\(\\*\\) AS users, count\\(DISTINCT user_id\\) AS links FROM links").WithoutArgs().WillReturnError(errors.New("err"))
			}

			got, err := tt.store.GetStat()
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.GetStat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.GetStat() = %v, want %v", got, tt.want)
			}
		})
	}
}
