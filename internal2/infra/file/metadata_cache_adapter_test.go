package file_test

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/file"
	infrafile "github.com/michaeldcanady/go-onedrive/internal2/infra/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMetadataCacheAdapter_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		cacheGetErr    error
		cacheGetVal    any
		expectOK       bool
		expectMetadata *file.Metadata
	}{
		{
			name:        "hit: valid metadata",
			path:        "/foo/meta.json",
			cacheGetErr: nil,
			cacheGetVal: file.Metadata{Path: "/foo/meta.json", ETag: "etag123"},
			expectOK:    true,
			expectMetadata: &file.Metadata{
				Path: "/foo/meta.json",
				ETag: "etag123",
			},
		},
		{
			name:           "miss: cache error",
			path:           "/foo/bad.json",
			cacheGetErr:    errors.New("not found"),
			cacheGetVal:    file.Metadata{},
			expectOK:       false,
			expectMetadata: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Metadata])

			mockCache.
				On("Get", mock.Anything, tt.path).
				Return(tt.cacheGetVal.(file.Metadata), tt.cacheGetErr)

			adapter := infrafile.NewMetadataCacheAdapter(mockCache)

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
		input     *file.Metadata
		setErr    error
		expectErr bool
	}{
		{
			name: "success: valid metadata",
			input: &file.Metadata{
				Path: "/foo/meta.json",
				ETag: "etag123",
			},
			setErr:    nil,
			expectErr: false,
		},
		{
			name: "error: cache.Set fails",
			input: &file.Metadata{
				Path: "/foo/fail",
				ETag: "etag",
			},
			setErr:    errors.New("cache write failed"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Metadata])

			if tt.input != nil {
				mockCache.
					On("Set", mock.Anything, tt.input.Path, *tt.input).
					Return(tt.setErr)
			}

			adapter := infrafile.NewMetadataCacheAdapter(mockCache)

			err := adapter.Put(context.Background(), tt.input.Path, tt.input)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[file.Metadata])

			mockCache.
				On("Delete", mock.Anything, tt.path).
				Return(tt.deleteErr)

			adapter := infrafile.NewMetadataCacheAdapter(mockCache)

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
