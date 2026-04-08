package cache

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockKeyValueStore is a mock implementation of KeyValueStore.
type MockKeyValueStore struct {
	mock.Mock
}

func (m *MockKeyValueStore) Get(ctx context.Context, key []byte) ([]byte, error) {
	args := m.Called(ctx, key)
	return args.Get(0).([]byte), args.Error(1)
}

func (m *MockKeyValueStore) Set(ctx context.Context, key []byte, value []byte) error {
	args := m.Called(ctx, key, value)
	return args.Error(0)
}

func (m *MockKeyValueStore) Delete(ctx context.Context, key []byte) error {
	args := m.Called(ctx, key)
	return args.Error(0)
}

func (m *MockKeyValueStore) List(ctx context.Context) ([][]byte, [][]byte, error) {
	args := m.Called(ctx)
	return args.Get(0).([][]byte), args.Get(1).([][]byte), args.Error(2)
}

func TestStore_Get(t *testing.T) {
	tests := []struct {
		name              string
		key               []byte
		storeResult       []byte
		storeErr          error
		expectedErr       string
		cancelCtx         bool
		keySerializer     SerializerFunc
		valueDeserializer func([]byte) error
	}{
		{
			name:        "Successful Get",
			key:         []byte("key1"),
			storeResult: []byte("value1"),
			keySerializer: func() ([]byte, error) {
				return []byte("key1"), nil
			},
			valueDeserializer: func(b []byte) error {
				assert.Equal(t, []byte("value1"), b)
				return nil
			},
		},
		{
			name: "Key Serializer Error",
			keySerializer: func() ([]byte, error) {
				return nil, errors.New("serialize error")
			},
			expectedErr: "serialize error",
		},
		{
			name: "Empty Key Error",
			keySerializer: func() ([]byte, error) {
				return []byte(""), nil
			},
			expectedErr: "key is empty",
		},
		{
			name: "Store Get Error",
			key:  []byte("key1"),
			keySerializer: func() ([]byte, error) {
				return []byte("key1"), nil
			},
			storeErr:    errors.New("store error"),
			expectedErr: "store error",
		},
		{
			name:      "Context Canceled",
			cancelCtx: true,
			keySerializer: func() ([]byte, error) {
				return []byte("key1"), nil
			},
			expectedErr: "context canceled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockKeyValueStore)
			s := NewStore(mockStore)
			ctx := context.Background()
			if tt.cancelCtx {
				var cancel context.CancelFunc
				ctx, cancel = context.WithCancel(ctx)
				cancel()
			}

			if tt.key != nil && !tt.cancelCtx && tt.expectedErr == "" {
				mockStore.On("Get", ctx, tt.key).Return(tt.storeResult, tt.storeErr)
			} else if tt.storeErr != nil {
				mockStore.On("Get", ctx, tt.key).Return([]byte(nil), tt.storeErr)
			}

			err := s.Get(ctx, tt.keySerializer, tt.valueDeserializer)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestStore_Set(t *testing.T) {
	tests := []struct {
		name            string
		key             []byte
		value           []byte
		storeErr        error
		expectedErr     string
		keySerializer   SerializerFunc
		valueSerializer SerializerFunc
	}{
		{
			name:  "Successful Set",
			key:   []byte("key1"),
			value: []byte("value1"),
			keySerializer: func() ([]byte, error) {
				return []byte("key1"), nil
			},
			valueSerializer: func() ([]byte, error) {
				return []byte("value1"), nil
			},
		},
		{
			name: "Store Set Error",
			key:  []byte("key1"),
			value: []byte("value1"),
			keySerializer: func() ([]byte, error) {
				return []byte("key1"), nil
			},
			valueSerializer: func() ([]byte, error) {
				return []byte("value1"), nil
			},
			storeErr:    errors.New("store error"),
			expectedErr: "store error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStore := new(MockKeyValueStore)
			s := NewStore(mockStore)
			ctx := context.Background()

			if tt.expectedErr == "" || tt.storeErr != nil {
				mockStore.On("Set", ctx, tt.key, tt.value).Return(tt.storeErr)
			}

			err := s.Set(ctx, tt.keySerializer, tt.valueSerializer)

			if tt.expectedErr != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
			mockStore.AssertExpectations(t)
		})
	}
}

func TestStore_Delete(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	s := NewStore(mockStore)
	ctx := context.Background()
	key := []byte("key1")

	mockStore.On("Delete", ctx, key).Return(nil)

	err := s.Delete(ctx, func() ([]byte, error) {
		return key, nil
	})

	assert.NoError(t, err)
	mockStore.AssertExpectations(t)
}

func TestStore_List(t *testing.T) {
	mockStore := new(MockKeyValueStore)
	s := NewStore(mockStore)
	ctx := context.Background()

	keys := [][]byte{[]byte("k1"), []byte("k2")}
	values := [][]byte{[]byte("v1"), []byte("v2")}

	mockStore.On("List", ctx).Return(keys, values, nil)

	var resultKeys [][]byte
	var resultValues [][]byte
	err := s.List(ctx, func(k, v []byte) error {
		resultKeys = append(resultKeys, k)
		resultValues = append(resultValues, v)
		return nil
	})

	assert.NoError(t, err)
	assert.Equal(t, keys, resultKeys)
	assert.Equal(t, values, resultValues)
	mockStore.AssertExpectations(t)
}
