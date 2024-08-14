package store

import "errors"

type InMemoryStore struct {
	data map[string]*StoreArgs
}

func (store *InMemoryStore) GetLink(arg *StoreArgs) (*StoreArgs, error) {
	link, ok := store.data[arg.ShortLink]
	if !ok {
		return nil, errors.New("key error")
	}

	arg.OriginalLink = link.OriginalLink

	return arg, nil

}

func (store *InMemoryStore) SetLink(arg *StoreArgs) (*StoreArgs, error) {
	store.data[arg.ShortLink] = arg

	return arg, nil
}

func (store *InMemoryStore) Ping() error {
	return nil
}

func (store *InMemoryStore) SetBatchLink(arg []StoreArgs) ([]*StoreArgs, error) {
	var res []*StoreArgs
	for _, val := range arg {
		store.data[val.ShortLink] = &val
		res = append(res, &val)
	}

	return res, nil
}

func CreateInMemoryStore() Store {
	store := &InMemoryStore{
		data: make(map[string]*StoreArgs),
	}

	return store
}
