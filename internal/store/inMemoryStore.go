package store

import "errors"

type InMemoryStore struct {
	data map[string]*StoreArgs
}

func (store *InMemoryStore) GetLink(arg *StoreArgs) (string, error) {
	link, ok := store.data[arg.ShortLink]
	if !ok {
		return "", errors.New("key error")
	}

	return link.OriginalLink, nil

}

func (store *InMemoryStore) SetLink(arg *StoreArgs) error {
	store.data[arg.ShortLink] = arg

	return nil
}

func (store *InMemoryStore) Ping() error {
	return nil
}

func (store *InMemoryStore) SetBatchLink(arg []*StoreArgs) ([]*StoreArgs, error) {
	for _, val := range arg {
		store.data[val.ShortLink] = val
	}

	return arg, nil
}

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]*StoreArgs),
	}

	return store
}
