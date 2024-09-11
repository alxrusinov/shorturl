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

func (store *InMemoryStore) GetLinks(userID string) ([]StoreRecord, error) {
	var result []StoreRecord

	for _, val := range store.data {
		if val.UUID == userID {
			result = append(result, *val)

		}

	}

	return result, nil
}

func (store *InMemoryStore) DeleteLinks(shorts [][]StoreRecord) error {
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

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]*StoreRecord),
	}

	return store
}
