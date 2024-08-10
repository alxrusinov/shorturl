package store

import "github.com/alxrusinov/shorturl/internal/config"

type Store interface {
	GetLink(key string) (string, error)
	SetLink(key string, link string) error
	Ping() error
}

func CreateStore(config *config.Config) Store {
	if config.DBPath != "" {
		return CreateDBStore(config.DBPath)
	}

	if config.FileStoragePath != "" {
		return CreateStore(config)
	}

	return CreateInMemoryStore()
}
