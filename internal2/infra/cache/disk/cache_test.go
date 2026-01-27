package disk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	mocks "github.com/michaeldcanady/go-onedrive/internal2/infra/cache/internal/mock"
)

func newMockCache(t *testing.T) (*Cache[string, string], *mocks.SerializerMock[string], *mocks.SerializerMock[string]) {
	keySer := &mocks.SerializerMock[string]{}
	valSer := &mocks.SerializerMock[string]{}

	path := t.TempDir() + "/cache.db"
	c, err := New(path, keySer, valSer)
	require.NoError(t, err)

	return c, keySer, valSer
}

func TestSetAndGetEntry_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		setupMock func(keySer, valSer *mocks.SerializerMock[string], key, value string)
	}{
		{
			name:  "simple write/read",
			key:   "foo",
			value: "bar",
			setupMock: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil).Maybe()
				valSer.On("Serialize", value).Return([]byte(value), nil)
				valSer.On("Deserialize", []byte(value)).Return(value, nil)
			},
		},
		{
			name:  "another key/value pair",
			key:   "hello",
			value: "world",
			setupMock: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil).Maybe()
				valSer.On("Serialize", value).Return([]byte(value), nil)
				valSer.On("Deserialize", []byte(value)).Return(value, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, keySer, valSer := newMockCache(t)
			ctx := context.Background()

			tt.setupMock(keySer, valSer, tt.key, tt.value)

			err := c.SetEntry(ctx, abstractions.NewEntry(tt.key, tt.value))
			require.NoError(t, err)

			entry, err := c.GetEntry(ctx, tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.value, entry.GetValue())

			keySer.AssertExpectations(t)
			valSer.AssertExpectations(t)
		})
	}
}

func TestWriteMultipleAndReadMultiple(t *testing.T) {
	tests := []struct {
		name  string
		pairs map[string]string
	}{
		{
			name: "three simple entries",
			pairs: map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		{
			name: "five mixed entries",
			pairs: map[string]string{
				"foo":   "bar",
				"hello": "world",
				"x":     "y",
				"key":   "value",
				"alpha": "beta",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, keySer, valSer := newMockCache(t)
			ctx := context.Background()

			// Set up mocks for all writes
			for k, v := range tt.pairs {
				keySer.On("Serialize", k).Return([]byte(k), nil).Maybe()
				valSer.On("Serialize", v).Return([]byte(v), nil)
			}

			// Set up mocks for all reads
			for k, v := range tt.pairs {
				keySer.On("Serialize", k).Return([]byte(k), nil).Maybe()
				valSer.On("Deserialize", []byte(v)).Return(v, nil)
			}

			// Write all entries
			for k, v := range tt.pairs {
				err := c.SetEntry(ctx, abstractions.NewEntry(k, v))
				require.NoError(t, err)
			}

			// Read all entries back
			for k, expected := range tt.pairs {
				entry, err := c.GetEntry(ctx, k)
				require.NoError(t, err)
				assert.Equal(t, expected, entry.GetValue())
			}

			keySer.AssertExpectations(t)
			valSer.AssertExpectations(t)
		})
	}
}

func TestRemove_TableDriven(t *testing.T) {
	tests := []struct {
		name      string
		key       string
		value     string
		setupMock func(keySer, valSer *mocks.SerializerMock[string], key, value string)
	}{
		{
			name:  "remove existing key",
			key:   "a",
			value: "1",
			setupMock: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil).Maybe()
				keySer.On("Deserialize", []byte(key)).Return(key, nil).Maybe()
				valSer.On("Serialize", value).Return([]byte(value), nil)
				valSer.On("Deserialize", []byte(value)).Return(value, nil).Maybe()
			},
		},
		{
			name:  "remove another key",
			key:   "b",
			value: "2",
			setupMock: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil).Maybe()
				keySer.On("Deserialize", []byte(key)).Return(key, nil).Maybe()
				valSer.On("Serialize", value).Return([]byte(value), nil)
				valSer.On("Deserialize", []byte(value)).Return(value, nil).Maybe()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, keySer, valSer := newMockCache(t)
			ctx := context.Background()

			tt.setupMock(keySer, valSer, tt.key, tt.value)

			require.NoError(t, c.SetEntry(ctx, abstractions.NewEntry(tt.key, tt.value)))

			require.NoError(t, c.Remove(tt.key))

			_, err := c.GetEntry(ctx, tt.key)
			assert.Error(t, err)

			keySer.AssertExpectations(t)
			valSer.AssertExpectations(t)
		})
	}
}

func TestPersistence_TableDriven(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		value  string
		setup1 func(keySer, valSer *mocks.SerializerMock[string], key, value string)
		setup2 func(keySer, valSer *mocks.SerializerMock[string], key, value string)
	}{
		{
			name:  "persist single key",
			key:   "alpha",
			value: "beta",
			setup1: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil)
				valSer.On("Serialize", value).Return([]byte(value), nil)
			},
			setup2: func(keySer, valSer *mocks.SerializerMock[string], key, value string) {
				keySer.On("Serialize", key).Return([]byte(key), nil)
				valSer.On("Deserialize", []byte(value)).Return(value, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := t.TempDir() + "/cache.db"
			ctx := context.Background()

			// First instance
			keySer1 := &mocks.SerializerMock[string]{}
			valSer1 := &mocks.SerializerMock[string]{}
			tt.setup1(keySer1, valSer1, tt.key, tt.value)

			c1, err := New(path, keySer1, valSer1)
			require.NoError(t, err)

			require.NoError(t, c1.SetEntry(ctx, abstractions.NewEntry(tt.key, tt.value)))

			// Reopen
			keySer2 := &mocks.SerializerMock[string]{}
			valSer2 := &mocks.SerializerMock[string]{}
			tt.setup2(keySer2, valSer2, tt.key, tt.value)

			c2, err := New(path, keySer2, valSer2)
			require.NoError(t, err)

			entry, err := c2.GetEntry(ctx, tt.key)
			require.NoError(t, err)
			assert.Equal(t, tt.value, entry.GetValue())

			keySer1.AssertExpectations(t)
			valSer1.AssertExpectations(t)
			keySer2.AssertExpectations(t)
			valSer2.AssertExpectations(t)
		})
	}
}
