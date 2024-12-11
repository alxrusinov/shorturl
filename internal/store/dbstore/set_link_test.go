package dbstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestDBStore_SetLinkSuccess(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
			args: args{
				arg: &model.StoreRecord{
					UUID:          "1",
					CorrelationID: "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					Deleted:       false,
				},
			},
			want: &model.StoreRecord{
				UUID:          "1",
				CorrelationID: "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				Deleted:       false,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO links").
				WithArgs(tt.args.arg.ShortLink, tt.args.arg.OriginalLink, tt.args.arg.CorrelationID, tt.args.arg.UUID).WillReturnResult(sqlmock.NewResult(1, 1))
			mock.ExpectQuery("SELECT short FROM links WHERE original = \\$1;").WithArgs(tt.args.arg.OriginalLink).WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("short"))

			got, err := tt.store.SetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetLinkError(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
			args: args{
				arg: &model.StoreRecord{
					UUID:          "1",
					CorrelationID: "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					Deleted:       false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO links").
				WithArgs(tt.args.arg.ShortLink, tt.args.arg.OriginalLink, tt.args.arg.CorrelationID, tt.args.arg.UUID).WillReturnError(errors.New("error"))

			got, err := tt.store.SetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetLinkDuplicate(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
			args: args{
				arg: &model.StoreRecord{
					UUID:          "1",
					CorrelationID: "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					Deleted:       false,
				},
			},
			want: &model.StoreRecord{
				UUID:          "1",
				CorrelationID: "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				Deleted:       false,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO links").
				WithArgs(tt.args.arg.ShortLink, tt.args.arg.OriginalLink, tt.args.arg.CorrelationID, tt.args.arg.UUID).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

			mock.ExpectQuery("SELECT short FROM links WHERE").WithArgs(tt.args.arg.OriginalLink).WillReturnRows(sqlmock.NewRows([]string{"short"}).AddRow("short"))

			got, err := tt.store.SetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetLinkSelectError(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
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
			args: args{
				arg: &model.StoreRecord{
					UUID:          "1",
					CorrelationID: "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					Deleted:       false,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectExec("INSERT INTO links").
				WithArgs(tt.args.arg.ShortLink, tt.args.arg.OriginalLink, tt.args.arg.CorrelationID, tt.args.arg.UUID).WillReturnError(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

			mock.ExpectQuery("SELECT short FROM links WHERE").WithArgs(tt.args.arg.OriginalLink).WillReturnError(errors.New("error"))

			got, err := tt.store.SetLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
