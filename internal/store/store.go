package store

import (
	"fmt"

	"github.com/alxrusinov/shorturl/internal/config"
)

type Store interface {
	GetLink(arg *StoreArgs) (string, error)
	SetLink(arg *StoreArgs) error
	Ping() error
}

type StoreArgs struct {
	ShortLink     string
	OriginalLink  string
	CorrelationId string
}

func CreateStore(config *config.Config) Store {
	fmt.Printf("%#v", config)
	if config.DBPath != "" {
		return CreateDBStore(config.DBPath)
	}

	if config.FileStoragePath != "" {
		return CreateFileStore(config.FileStoragePath)
	}

	return CreateInMemoryStore()
}
