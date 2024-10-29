package mockstore

import (
	"github.com/alxrusinov/shorturl/internal/model"
	"github.com/stretchr/testify/mock"
)

// MockStore - store mockking real store
type MockStore struct {
	mock.Mock
}

// GetLink returns link from store
func (ms *MockStore) GetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(0).(*model.StoreRecord), args.Error(1)
}

// SetLink adds link to store
func (ms *MockStore) SetLink(arg *model.StoreRecord) (*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(0).(*model.StoreRecord), args.Error(1)
}

// SetBatchLink adds links to store by batch
func (ms *MockStore) SetBatchLink(arg []*model.StoreRecord) ([]*model.StoreRecord, error) {
	args := ms.Called(arg)

	return args.Get(0).([]*model.StoreRecord), args.Error(1)
}

// Ping pings store
func (ms *MockStore) Ping() error {
	args := ms.Called()

	return args.Error(0)
}

// GetLinks returns all users links
func (ms *MockStore) GetLinks(userID string) ([]model.StoreRecord, error) {
	args := ms.Called(userID)

	return args.Get(0).([]model.StoreRecord), args.Error(1)
}

// DeleteLinks deletes links by batch
func (ms *MockStore) DeleteLinks(shorts [][]model.StoreRecord) error {
	args := ms.Called(shorts)

	return args.Error(0)
}

// NewMockStore returns new mockestore instnce
func NewMockStore() *MockStore {
	return new(MockStore)
}
