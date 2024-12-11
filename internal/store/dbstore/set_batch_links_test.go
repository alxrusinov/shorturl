package dbstore

import (
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alxrusinov/shorturl/internal/model"
)

func TestDBStore_SetBatchLinkSuccess(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		arg []*model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    []*model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{
				arg: []*model.StoreRecord{{
					UUID:          "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					CorrelationID: "1",
					Deleted:       false,
				}},
			},
			want: []*model.StoreRecord{{
				UUID:          "1",
				ShortLink:     "short",
				OriginalLink:  "original",
				CorrelationID: "1",
				Deleted:       false,
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT").WithArgs("short", "original", "1", "1").WillReturnRows(sqlmock.NewRows([]string{"short", "original", "correlation_id", "user_id"}).AddRow("short", "original", "1", "1"))
			mock.ExpectCommit()
			got, err := tt.store.SetBatchLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetBatchLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetBatchLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetBatchLinkBeginFail(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		arg []*model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    []*model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{
				arg: []*model.StoreRecord{{
					UUID:          "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					CorrelationID: "1",
					Deleted:       false,
				}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin().WillReturnError(errors.New("err"))
			mock.ExpectQuery("INSERT").WithArgs("short", "original", "1", "1").WillReturnRows(sqlmock.NewRows([]string{"short", "original", "correlation_id", "user_id"}).AddRow("short", "original", "1", "1"))
			mock.ExpectRollback()
			got, err := tt.store.SetBatchLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetBatchLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetBatchLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetBatchLinkCommitFail(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		arg []*model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    []*model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{
				arg: []*model.StoreRecord{{
					UUID:          "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					CorrelationID: "1",
					Deleted:       false,
				}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT").WithArgs("short", "original", "1", "1").WillReturnRows(sqlmock.NewRows([]string{"short", "original", "correlation_id", "user_id"}).AddRow("short", "original", "1", "1"))
			mock.ExpectCommit().WillReturnError(errors.New("err"))
			got, err := tt.store.SetBatchLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetBatchLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetBatchLink() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDBStore_SetBatchLinkQueryErr(t *testing.T) {
	teststore, mock, err := NewDBStoreMock()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer teststore.db.Close()
	type args struct {
		arg []*model.StoreRecord
	}
	tests := []struct {
		name    string
		store   *DBStore
		args    args
		want    []*model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: teststore,
			args: args{
				arg: []*model.StoreRecord{{
					UUID:          "1",
					ShortLink:     "short",
					OriginalLink:  "original",
					CorrelationID: "1",
					Deleted:       false,
				}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock.ExpectBegin()
			mock.ExpectQuery("INSERT").WithArgs("short", "original", "1", "1").WillReturnError(errors.New("err"))
			mock.ExpectCommit()
			got, err := tt.store.SetBatchLink(tt.args.arg)
			if (err != nil) != tt.wantErr {
				t.Errorf("DBStore.SetBatchLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DBStore.SetBatchLink() = %v, want %v", got, tt.want)
			}
		})
	}
}
