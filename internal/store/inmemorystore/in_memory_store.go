package inmemorystore

import (
	"errors"

	"github.com/alxrusinov/shorturl/internal/model"
)

type InMemoryStore struct {
	data map[string]*model.StoreRecord
}

func (store *InMemoryStore) GetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	link, ok := store.data[arg.ShortLink]
	if !ok {
		return nil, errors.New("key error")
	}

	arg.OriginalLink = link.OriginalLink

	return arg, nil

}

func (store *InMemoryStore) SetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	store.data[arg.ShortLink] = arg

	return arg, nil
}

func (store *InMemoryStore) Ping() error {
	return nil
}

func (store *InMemoryStore) SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error) {
	for _, val := range arg {
		store.data[val.ShortLink] = val
	}

	return arg, nil
}

func (store *InMemoryStore) GetLinks(userID string) ([]model.StoreRecord, error) {
	var result []model.StoreRecord

	for _, val := range store.data {
		if val.UUID == userID {
			result = append(result, *val)

		}

	}

	return result, nil
}

func (store *InMemoryStore) DeleteLinks(shorts [][]model.StoreRecord) error {
	for _, val := range shorts {
		for _, short := range val {
			if record, ok := store.data[short.ShortLink]; ok {
				if record.UUID == short.UUID && record.ShortLink == short.ShortLink {
					store.data[short.ShortLink].Deleted = true
				}
			} else {
				return errors.New("key error")
			}
		}
	}

	return nil
}

func NewInMemoryStore() *InMemoryStore {
	store := &InMemoryStore{
		data: make(map[string]*model.StoreRecord),
	}

	return store
}
