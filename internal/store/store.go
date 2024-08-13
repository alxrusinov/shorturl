package store

import (
	"github.com/alxrusinov/shorturl/internal/config"
)

type Store interface {
	GetLink(arg *StoreArgs) (*StoreArgs, error)
	SetLink(arg *StoreArgs) (*StoreArgs, error)
	SetBatchLink(arg []*StoreArgs) ([]*StoreArgs, error)
	Ping() error
}

type StoreArgs struct {
	ShortLink     string `json:"short_url,omitempty"`
	OriginalLink  string `json:"original_url,omitempty"`
	CorrelationID string `json:"correlation_id"`
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
