package file_test

import (
	"context"
	"errors"
	"testing"

	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCache2 struct {
	mock.Mock
}

func (m *MockCache2) Get(ctx context.Context, keyFn abstractions.Serializer2, unmarshalFn abstractions.Deserializer2) error {
	args := m.Called(ctx, keyFn, unmarshalFn)
	return args.Error(0)
}

func (m *MockCache2) Set(ctx context.Context, keyFn abstractions.Serializer2, marshalFn abstractions.Serializer2) error {
	args := m.Called(ctx, keyFn, marshalFn)
	return args.Error(0)
}

func (m *MockCache2) Delete(ctx context.Context, keyFn abstractions.Serializer2) error {
	args := m.Called(ctx, keyFn)
	return args.Error(0)
}

func TestContentsCacheAdapter_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		cacheGetErr    error
		cacheGetBytes  []byte
		expectOK       bool
		expectContents *domainfile.Contents
	}{
		{
			name:          "hit: valid JSON",
			path:          "/foo/bar.txt",
			cacheGetErr:   nil,
			cacheGetBytes: mustJSON(t, domainfile.Contents{CTag: "etag123", Data: []byte("hello")}),
			expectOK:      true,
			expectContents: &domainfile.Contents{
				CTag: "etag123",
				Data: []byte("hello"),
			},
		},
		{
			name:           "miss: cache error",
			path:           "/foo/missing.txt",
			cacheGetErr:    errors.New("not found"),
			cacheGetBytes:  nil,
			expectOK:       false,
			expectContents: nil,
		},
		{
			name:           "miss: invalid JSON",
			path:           "/foo/bad.json",
			cacheGetErr:    errors.New("invalid json"),
			cacheGetBytes:  []byte("{not valid json"),
			expectOK:       false,
			expectContents: nil,
		},
		{
			name:           "hit: empty struct",
			path:           "/foo/empty",
			cacheGetErr:    nil,
			cacheGetBytes:  mustJSON(t, domainfile.Contents{}),
			expectOK:       true,
			expectContents: &domainfile.Contents{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2)

			mockCache.
				On("Get", mock.Anything, mock.AnythingOfType("abstractions.Serializer2"), mock.AnythingOfType("abstractions.Deserializer2")).
				Run(func(args mock.Arguments) {
					unmarshalFn := args.Get(2).(abstractions.Deserializer2)
					_ = unmarshalFn(tt.cacheGetBytes)
				}).
				Return(tt.cacheGetErr)

			adapter := file.NewContentsCacheAdapter(mockCache)

			got, ok := adapter.Get(context.Background(), tt.path)

			require.Equal(t, tt.expectOK, ok)

			if !tt.expectOK {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			assert.Equal(t, tt.expectContents.CTag, got.CTag)
			assert.Equal(t, tt.expectContents.Data, got.Data)

			mockCache.AssertExpectations(t)
		})
	}
}

func TestContentsCacheAdapter_Put(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		input     *domainfile.Contents
		setErr    error
		expectErr bool
	}{
		{
			name: "success: valid contents",
			path: "/foo/bar.txt",
			input: &domainfile.Contents{
				CTag: "etag123",
				Data: []byte("hello"),
			},
			setErr:    nil,
			expectErr: false,
		},
		{
			name: "error: cache.Set fails",
			path: "/foo/fail",
			input: &domainfile.Contents{
				CTag: "etag",
				Data: []byte("data"),
			},
			setErr:    errors.New("cache write failed"),
			expectErr: true,
		},
		{
			name:      "success: empty contents",
			path:      "/foo/empty",
			input:     &domainfile.Contents{},
			setErr:    nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2)

			mockCache.
				On("Set", mock.Anything, mock.AnythingOfType("abstractions.Serializer2"), mock.AnythingOfType("abstractions.Serializer2")).
				Run(func(args mock.Arguments) {
					marshalFn := args.Get(2).(abstractions.Serializer2)
					_, _ = marshalFn()
				}).
				Return(tt.setErr)

			adapter := file.NewContentsCacheAdapter(mockCache)

			err := adapter.Put(context.Background(), tt.path, tt.input)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockCache.AssertExpectations(t)
		})
	}
}

func TestContentsCacheAdapter_Invalidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		deleteErr error
		expectErr bool
	}{
		{
			name:      "success: delete works",
			path:      "/foo/bar.txt",
			deleteErr: nil,
			expectErr: false,
		},
		{
			name:      "error: delete fails",
			path:      "/foo/fail",
			deleteErr: errors.New("delete error"),
			expectErr: true,
		},
		{
			name:      "success: deleting non-existent key returns nil",
			path:      "/foo/missing",
			deleteErr: nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2)

			mockCache.
				On("Delete", mock.Anything, mock.AnythingOfType("abstractions.Serializer2")).
				Return(tt.deleteErr)

			adapter := file.NewContentsCacheAdapter(mockCache)

			err := adapter.Invalidate(context.Background(), tt.path)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockCache.AssertExpectations(t)
		})
	}
}
