package store

import "errors"

type InMemoryStore struct {
	data map[string]*StoreRecord
}

func (store *InMemoryStore) GetLink(arg *StoreRecord) (*StoreRecord, error) {
	link, ok := store.data[arg.ShortLink]
	if !ok {
		return nil, errors.New("key error")
	}

	arg.OriginalLink = link.OriginalLink

	return arg, nil

}

func (store *InMemoryStore) SetLink(arg *StoreRecord) (*StoreRecord, error) {
	store.data[arg.ShortLink] = arg

	return arg, nil
}

func (store *InMemoryStore) Ping() error {
	return nil
}

func (store *InMemoryStore) SetBatchLink(arg []*StoreRecord) ([]*StoreRecord, error) {
	for _, val := range arg {
		store.data[val.ShortLink] = val
	}

	return arg, nil
}

func (store *InMemoryStore) GetLinks(userId string) ([]StoreRecord, error) {
	var result []StoreRecord

	for _, val := range store.data {
		if val.UUID == userId {
			result = append(result, *val)

		}

	}

	return result, nil
}

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]*StoreRecord),
	}

	return store
}
