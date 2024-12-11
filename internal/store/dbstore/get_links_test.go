package dbstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
)

func TestDBStore_GetLinks(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    []model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{
				userID: "1",
			},
			want: []model.StoreRecord{
				{
					UUID:          "1",
					OriginalLink:  "original",
					ShortLink:     "short",
					CorrelationID: "1",
					Deleted:       false,
				},
			},
			wantErr: false,
		},
		{
			name:  "2# error",
			store: teststore,
			args: args{
				userID: "1",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:  "3# scan error",
			store: teststore,
			args: args{
				userID: "1",
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectQuery("SELECT user_id, short, original, correlation_id, is_deleted FROM links WHERE user_id = \\$1").WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"user_id", "short", "original", "correlation_id", "is_deleted"}).AddRow("1", "short", "original", "1", false))
			} else {
				if tt.name == tests[1].name {
					mock.ExpectQuery("SELECT user_id, short, original, correlation_id, is_deleted FROM links WHERE user_id = \\$1").WithArgs("1").WillReturnError(errors.New("error"))
				}

				if tt.name == tests[2].name {
					mock.ExpectQuery("SELECT user_id, short, original, correlation_id, is_deleted FROM links WHERE user_id = \\$1").WithArgs("1").WillReturnRows(sqlmock.NewRows([]string{"user_id", "short", "original", "correlation_id", "is_deleted"}).AddRow("1", "short", "original", "1", "222"))
				}

			}

			got, err := tt.store.GetLinks(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.GetLinks() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.GetLinks() = %v, want %v", got, tt.want)
			}
		})
	}
}
