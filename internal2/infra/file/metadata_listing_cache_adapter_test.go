package file_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/infra/cache/abstractions"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMetadataListCacheAdapter_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		path          string
		cacheGetBytes []byte
		expectOK      bool
		expectListing *file.Listing
	}{
		{
			name:          "hit: valid JSON",
			path:          "/foo/listing",
			cacheGetBytes: mustJSON(t, file.Listing{ETag: "etag123", ChildIDs: []string{"a", "b"}}),
			expectOK:      true,
			expectListing: &file.Listing{
				ETag:     "etag123",
				ChildIDs: []string{"a", "b"},
			},
		},
		{
			name:          "miss: invalid JSON",
			path:          "/foo/bad",
			cacheGetBytes: []byte("{not valid json"),
			expectOK:      false,
			expectListing: nil,
		},
		{
			name:          "hit: empty listing",
			path:          "/foo/empty",
			cacheGetBytes: mustJSON(t, file.Listing{}),
			expectOK:      true,
			expectListing: &file.Listing{},
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
					err := unmarshalFn(tt.cacheGetBytes)

					// If unmarshal fails, override return value
					if err != nil {
						mockCache.ExpectedCalls[0].ReturnArguments = mock.Arguments{err}
					}
				}).
				Return(nil)

			adapter := file.NewMetadataListCacheAdapter(mockCache)

			got, ok := adapter.Get(context.Background(), tt.path)

			require.Equal(t, tt.expectOK, ok)

			if !tt.expectOK {
				require.Nil(t, got)
				return
			}

			require.NotNil(t, got)
			assert.Equal(t, tt.expectListing.ETag, got.ETag)
			assert.Equal(t, tt.expectListing.ChildIDs, got.ChildIDs)

			mockCache.AssertExpectations(t)
		})
	}
}

func TestMetadataListCacheAdapter_Put(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		listing   *file.Listing
		setErr    error
		expectErr bool
	}{
		{
			name: "success: valid listing",
			path: "/foo/listing",
			listing: &file.Listing{
				ETag:     "etag123",
				ChildIDs: []string{"a", "b"},
			},
			setErr:    nil,
			expectErr: false,
		},
		{
			name: "error: cache.Set fails",
			path: "/foo/fail",
			listing: &file.Listing{
				ETag:     "etag",
				ChildIDs: []string{"x"},
			},
			setErr:    errors.New("cache write failed"),
			expectErr: true,
		},
		{
			name:      "success: empty listing",
			path:      "/foo/empty",
			listing:   &file.Listing{},
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

			adapter := file.NewMetadataListCacheAdapter(mockCache)

			err := adapter.Put(context.Background(), tt.path, tt.listing)

			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockCache.AssertExpectations(t)
		})
	}
}

func TestMetadataListCacheAdapter_Invalidate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		path      string
		deleteErr error
		expectErr bool
	}{
		{
			name:      "success: delete works",
			path:      "/foo/listing",
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

			adapter := file.NewMetadataListCacheAdapter(mockCache)

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
