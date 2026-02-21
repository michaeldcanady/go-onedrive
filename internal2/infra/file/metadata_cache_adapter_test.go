package file_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	domainfile "github.com/michaeldcanady/go-onedrive/internal2/domain/file"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func mustJSON(t *testing.T, v any) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

func TestMetadataCacheAdapter_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		cacheGetErr    error
		cacheGetBytes  []byte
		expectOK       bool
		expectMetadata *domainfile.Metadata
	}{
		{
			name:          "hit: valid JSON",
			path:          "/foo/meta.json",
			cacheGetErr:   nil,
			cacheGetBytes: mustJSON(t, domainfile.Metadata{Path: "/foo/meta.json", ETag: "etag123"}),
			expectOK:      true,
			expectMetadata: &domainfile.Metadata{
				Path: "/foo/meta.json",
				ETag: "etag123",
			},
		},
		{
			name:           "miss: invalid JSON",
			path:           "/foo/bad.json",
			cacheGetErr:    errors.New("invalid json"),
			cacheGetBytes:  []byte("{not valid json"),
			expectOK:       false,
			expectMetadata: nil,
		},
		{
			name:           "hit: empty struct",
			path:           "/foo/empty",
			cacheGetBytes:  mustJSON(t, domainfile.Metadata{}),
			expectOK:       true,
			expectMetadata: &domainfile.Metadata{},
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

			adapter := file.NewMetadataCacheAdapter(mockCache)

			got, ok := adapter.Get(context.Background(), tt.path)

			require.Equal(t, tt.expectOK, ok)

			if !tt.expectOK {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			assert.Equal(t, tt.expectMetadata.Path, got.Path)
			assert.Equal(t, tt.expectMetadata.ETag, got.ETag)

			mockCache.AssertExpectations(t)
		})
	}
}

func TestMetadataCacheAdapter_Put(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     *domainfile.Metadata
		setErr    error
		expectErr bool
	}{
		{
			name: "success: valid metadata",
			input: &domainfile.Metadata{
				Path: "/foo/meta.json",
				ETag: "etag123",
			},
			setErr:    nil,
			expectErr: false,
		},
		{
			name: "error: cache.Set fails",
			input: &domainfile.Metadata{
				Path: "/foo/fail",
				ETag: "etag",
			},
			setErr:    errors.New("cache write failed"),
			expectErr: true,
		},
		{
			name:      "success: empty metadata",
			input:     &domainfile.Metadata{},
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

			adapter := file.NewMetadataCacheAdapter(mockCache)

			err := adapter.Put(context.Background(), tt.input)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockCache.AssertExpectations(t)
		})
	}
}

func TestMetadataCacheAdapter_Invalidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		deleteErr error
		expectErr bool
	}{
		{
			name:      "success: delete works",
			path:      "/foo/meta.json",
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

			adapter := file.NewMetadataCacheAdapter(mockCache)

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
