package dbstore

import (
	"errors"
	"testing"
)

func TestDBStore_Ping(t *testing.T) {
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
			name:  "1# success",
			store: teststore,
		},
		{
			name:    "1# error",
			store:   teststore,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.wantErr {
				mock.ExpectPing().WillReturnError(errors.New("err"))
			} else {
				mock.ExpectPing().WillReturnError(nil)
			}

			if err := tt.store.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("DBStore.Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
