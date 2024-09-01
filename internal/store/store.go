package store

import (
	"github.com/alxrusinov/shorturl/internal/config"
)

type Store interface {
	GetLink(arg *StoreRecord) (*StoreRecord, error)
	SetLink(arg *StoreRecord) (*StoreRecord, error)
	SetBatchLink(arg []*StoreRecord) ([]*StoreRecord, error)
	Ping() error
	GetLinks(userID string) ([]StoreRecord, error)
}

type StoreRecord struct {
	UUID          string `json:"user_id" db:"user_id"`
	ShortLink     string `json:"short_url,omitempty" db:"short"`
	OriginalLink  string `json:"original_url,omitempty" db:"original"`
	CorrelationID string `json:"correlation_id" db:"correlation_id"`
	Deleted       bool   `json:"is_deleted" db:"is_deleted"`
}

func CreateStore(config *config.Config) Store {
	if config.DBPath != "" {
		return CreateDBStore(config.DBPath)
	}

	if config.FileStoragePath != "" {
		return CreateFileStore(config.FileStoragePath)
	}

	return CreateInMemoryStore()
}
