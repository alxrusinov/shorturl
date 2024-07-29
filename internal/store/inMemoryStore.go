package store

import "errors"

type InMemoryStore struct {
	data map[string]string
}

func (store *InMemoryStore) GetLink(key string) (string, error) {
	link, ok := store.data[key]
	if !ok {
		return "", errors.New("key error")
	}

	return link, nil

}

func (store *InMemoryStore) SetLink(key string, link string) {
	store.data[key] = link
}

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]string),
	}

	return store
}
