package filestore

import (
	"fmt"
	"os"
	"testing"

	"github.com/alxrusinov/shorturl/internal/customerrors"
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestFileStore_GetLink(t *testing.T) {
	const fileStorePath string = "fileStore"
	const wrongStorePath string = "foo"
	fileStore := NewFileStore(fileStorePath)
	wronFileStore := NewFileStore(wrongStorePath)

	file, _ := os.Create(fileStorePath)

	trueRecord := &model.StoreRecord{
		CorrelationID: "1",
		OriginalLink:  "http://example.com",
		ShortLink:     "short_1",
	}

	notFoundRecord := &model.StoreRecord{
		CorrelationID: "2",
		OriginalLink:  "http://example.com",
		ShortLink:     "short_2",
	}

	fileStore.SetLink(trueRecord)

	tests := []struct {
		name    string
		store   *FileStore
		args    *model.StoreRecord
		want    *model.StoreRecord
		wantErr bool
	}{
		{
			name:  "1# success",
			store: fileStore,
			args:  trueRecord,
			want:  trueRecord,
		},
		{
			name:    "2# not found",
			store:   fileStore,
			args:    notFoundRecord,
			want:    trueRecord,
			wantErr: true,
		},
		{
			name:    "2# err open file",
			store:   wronFileStore,
			args:    notFoundRecord,
			want:    trueRecord,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.store.GetLink(tt.args)
			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tt.want.OriginalLink, got.OriginalLink)
				assert.Equal(t, tt.want.ShortLink, got.ShortLink)
				assert.Equal(t, tt.want.CorrelationID, got.CorrelationID)
				assert.Equal(t, tt.want.UUID, got.UUID)
			}

		})
	}

	file.Close()

	err := os.Remove(fileStorePath)

	if err != nil {
		fmt.Printf("ERROR - %#v\n", err.Error())
	}
}

func TestFileStore_SetLink(t *testing.T) {
	const fileStorePath string = "fileStoreSet"

	fileStore := NewFileStore(fileStorePath)

	file, _ := os.Create(fileStorePath)

	trueRecord := &model.StoreRecord{
		CorrelationID: "1",
		OriginalLink:  "http://example.com",
		ShortLink:     "short_1",
	}

	tests := []struct {
		name             string
		store            *FileStore
		args             *model.StoreRecord
		want             *model.StoreRecord
		wantErr          bool
		wantDuplicateErr bool
	}{
		{
			name:  "1# success",
			store: fileStore,
			args:  trueRecord,
			want:  trueRecord,
		},
		{
			name:             "2# err",
			store:            fileStore,
			args:             trueRecord,
			want:             trueRecord,
			wantDuplicateErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duplicateErr := &customerrors.DuplicateValueError{}
			got, err := tt.store.SetLink(tt.args)
			switch {
			case tt.wantErr:
				assert.NotNil(t, err)
			case tt.wantDuplicateErr:
				assert.ErrorAs(t, err, &duplicateErr)
			default:
				assert.Nil(t, err)
				assert.Equal(t, tt.want.OriginalLink, got.OriginalLink)
				assert.Equal(t, tt.want.ShortLink, got.ShortLink)
				assert.Equal(t, tt.want.CorrelationID, got.CorrelationID)
				assert.Equal(t, tt.want.UUID, got.UUID)
			}

		})
	}

	if err := file.Close(); err == nil {
		err := os.Remove(fileStorePath)

		if err != nil {
			fmt.Printf("ERROR - %#v\n", err.Error())
		}
	}

}

func TestFileStore_Ping(t *testing.T) {
	const fileStorePath string = "fileStorePing"
	const wrongFileStorePath string = "bar"

	fileStore := NewFileStore(fileStorePath)
	wrongFileStore := NewFileStore(wrongFileStorePath)

	file, _ := os.Create(fileStorePath)

	tests := []struct {
		name    string
		store   *FileStore
		wantErr bool
	}{
		{
			name:    "1# success",
			store:   fileStore,
			wantErr: false,
		},
		{
			name:    "2# err",
			store:   wrongFileStore,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.store.Ping()

			if tt.wantErr {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

	if err := file.Close(); err == nil {
		err = os.Remove(fileStorePath)

		if err != nil {
			fmt.Printf("ERROR - %#v\n", err.Error())
		}
	}
}
