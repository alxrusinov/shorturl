package store

import (
	"errors"
)

type Store interface {
	GetLink(key string) (string, error)
	SetLink(key string, link string)
}

type Cache struct {
	data map[string]string
}

func (store *Cache) GetLink(key string) (string, error) {
link, ok := store.data[key]
	if !ok {
		return "", errors.New("key error")
	}

	return link, nil

}

func (store *Cache) SetLink(key string, link string) {
	store.data[key] = link
}

func CreateStore() Store {
	store := &Cache{
		data: make(map[string]string),
	}

	return store
}
