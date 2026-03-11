package infra

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockCache2[T any] struct {
	mock.Mock
}

func (m *MockCache2[T]) Get(ctx context.Context, key string) (T, error) {
	args := m.Called(ctx, key)
	return args.Get(0).(T), args.Error(1)
}

func (m *MockCache2[T]) Set(ctx context.Context, key string, value T) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockCache2[T]) Delete(ctx context.Context, key string) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockCache2[T]) List(ctx context.Context, callback func(key string, value T) error) error {
	args := m.Called(ctx, callback)
	return args.Error(0)
}
