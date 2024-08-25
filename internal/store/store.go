package store

import (
	"github.com/alxrusinov/shorturl/internal/config"
)

type Store interface {
	GetLink(arg *StoreRecord) (*StoreRecord, error)
	SetLink(arg *StoreRecord) (*StoreRecord, error)
	SetBatchLink(arg []*StoreRecord) ([]*StoreRecord, error)
	Ping() error
}

type StoreRecord struct {
	UUID          string `json:"user_id"`
	ShortLink     string `json:"short_url,omitempty"`
	OriginalLink  string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id"`
	Deleted       bool   `json:"is_deleted"`
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
