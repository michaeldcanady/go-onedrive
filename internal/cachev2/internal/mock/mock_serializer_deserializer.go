package mocks

import (
	"github.com/stretchr/testify/mock"
)

// MockSerializerDeserializer
type SerializerMock[T any] struct {
	mock.Mock
}

func (m *SerializerMock[T]) Serialize(v T) ([]byte, error) {
	args := m.Called(v)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *SerializerMock[T]) Deserialize(b []byte) (T, error) {
	args := m.Called(b)
	return args.Get(0).(T), args.Error(1)
}
