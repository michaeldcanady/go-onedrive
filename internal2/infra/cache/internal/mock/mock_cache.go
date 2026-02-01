package mocks

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/stretchr/testify/mock"
)

type MockCache[K comparable, V any] struct {
	mock.Mock
}

func (m *MockCache[K, V]) NewEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*abstractions.Entry[K, V]), args.Error(1)
}

func (m *MockCache[K, V]) GetEntry(ctx context.Context, key K) (*abstractions.Entry[K, V], error) {
	args := m.Called(ctx, key)
	return args.Get(0).(*abstractions.Entry[K, V]), args.Error(1)
}

func (m *MockCache[K, V]) SetEntry(ctx context.Context, e *abstractions.Entry[K, V]) error {
	args := m.Called(ctx, e)
	return args.Error(0)
}

func (m *MockCache[K, V]) Clear(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

func (m *MockCache[K, V]) Remove(key K) error {
	args := m.Called(key)
	return args.Error(0)
}

func (m *MockCache[K, V]) KeySerializer() abstractions.Serializer[K] {
	args := m.Called()
	return args.Get(0).(abstractions.SerializerDeserializer[K])
}
