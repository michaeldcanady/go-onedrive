package file_test

import (
	"context"
	"errors"
	"testing"

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
		cacheGetErr   error
		cacheGetVal   any
		expectOK      bool
		expectListing *file.Listing
	}{
		{
			name:        "hit: valid listing",
			path:        "/foo/listing",
			cacheGetErr: nil,
			cacheGetVal: file.Listing{ETag: "etag123", ChildIDs: []string{"a", "b"}},
			expectOK:    true,
			expectListing: &file.Listing{
				ETag:     "etag123",
				ChildIDs: []string{"a", "b"},
			},
		},
		{
			name:          "miss: cache error",
			path:          "/foo/bad",
			cacheGetErr:   errors.New("not found"),
			cacheGetVal:   file.Listing{},
			expectOK:      false,
			expectListing: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Listing])

			mockCache.
				On("Get", mock.Anything, tt.path).
				Return(tt.cacheGetVal.(file.Listing), tt.cacheGetErr)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Listing])

			if tt.listing != nil {
				mockCache.
					On("Set", mock.Anything, tt.path, *tt.listing).
					Return(tt.setErr)
			}

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Listing])

			mockCache.
				On("Delete", mock.Anything, tt.path).
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
