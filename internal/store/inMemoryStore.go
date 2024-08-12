package store

import "errors"

type InMemoryStore struct {
	data map[string]string
}

func (store *InMemoryStore) GetLink(arg *StoreArgs) (string, error) {
	link, ok := store.data[arg.ShortLink]
	if !ok {
		return "", errors.New("key error")
	}

	return link, nil

}

func (store *InMemoryStore) SetLink(arg *StoreArgs) error {
	store.data[arg.CorrelationId] = arg.OriginalLink

	return nil
}

func (store *InMemoryStore) Ping() error {
	return nil
}

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]string),
	}

	return store
}
