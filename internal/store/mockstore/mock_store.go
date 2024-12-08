package mockstore

import (
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/stretchr/testify/mock"
)

const (
	indexZero = iota
	indexOne
)

// MockStore - store mockking real store
type MockStore struct {
	mock.Mock
}

// GetLink returns link from store
func (ms *MockStore) GetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(indexZero).(*model.StoreRecord), args.Error(indexOne)
}

// SetLink adds link to store
func (ms *MockStore) SetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(indexZero).(*model.StoreRecord), args.Error(indexOne)
}

// SetBatchLink adds links to store by batch
func (ms *MockStore) SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(indexZero).([]*model.StoreRecord), args.Error(indexOne)
}

// Ping pings store
func (ms *MockStore) Ping() error {
	args := ms.Called()

	return args.Error(0)
}

// GetLinks returns all users links
func (ms *MockStore) GetLinks(userID string) ([]model.StoreRecord, error) {
	args := ms.Called(userID)

	return args.Get(indexZero).([]model.StoreRecord), args.Error(indexOne)
}

// DeleteLinks deletes links by batch
func (ms *MockStore) DeleteLinks(shorts [][]model.StoreRecord) error {
	args := ms.Called(shorts)

	return args.Error(indexZero)
}

// GetStat - gets dtatistics of urls and users
func (ms *MockStore) GetStat() (*model.StatResponse, error) {
	args := ms.Called()

	return args.Get(indexZero).(*model.StatResponse), args.Error(indexOne)
}

// NewMockStore returns new mockestore instnce
func NewMockStore() *MockStore {
	return new(MockStore)
}
