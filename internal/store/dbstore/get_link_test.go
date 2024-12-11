package dbstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
)

func TestDBStore_GetLink(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		arg *model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    *model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{arg: &model.StoreRecord{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}},
			want: &model.StoreRecord{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			},
		},
		{
			name:  "2# error",
			store: teststore,
			args: args{arg: &model.StoreRecord{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !tt.wantErr {
				mock.ExpectQuery("SELECT original FROM links").WillReturnRows(sqlmock.NewRows([]string{"original"}).AddRow(tt.want.OriginalLink))
			} else {
				mock.ExpectQuery("SELECT original FROM links").WillReturnError(errors.New("error"))
			}
			got, err := tt.store.GetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.GetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.GetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
