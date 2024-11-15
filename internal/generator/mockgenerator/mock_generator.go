package mockgenerator

import "github.com/stretchr/testify/mock"

const (
	indexZero = iota
	indexOne
)

// MochGenerator - mocked generator
type MockGenerator struct {
	mock.Mock
}

// GenerateRandomString - mocked method of generator
func (mg *MockGenerator) GenerateRandomString() (string, error) {

	args := mg.Called()

	return args.String(indexZero), args.Error(indexOne)
}

// GenerateUserID - mocked method of generator
func (mg *MockGenerator) GenerateUserID() (string, error) {
	args := mg.Called()

	return args.String(indexZero), args.Error(indexOne)
}

// NewMockGenerator return mocked generator
func NewMockGenerator() *MockGenerator {
	return &MockGenerator{}
}
