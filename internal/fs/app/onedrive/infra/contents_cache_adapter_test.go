package infra

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestContentsCacheAdapter_Get(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		path           string
		cacheGetErr    error
		cacheGetVal    any
		expectOK       bool
		expectContents *domain.Contents
	}{
		{
			name:        "hit: valid contents",
			path:        "/foo/bar.txt",
			cacheGetErr: nil,
			cacheGetVal: domain.Contents{CTag: "etag123", Data: []byte("hello")},
			expectOK:    true,
			expectContents: &domain.Contents{
				CTag: "etag123",
				Data: []byte("hello"),
			},
		},
		{
			name:           "miss: cache error",
			path:           "/foo/missing.txt",
			cacheGetErr:    errors.New("not found"),
			cacheGetVal:    domain.Contents{},
			expectOK:       false,
			expectContents: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[domain.Contents])

			mockCache.
				On("Get", mock.Anything, tt.path).
				Return(tt.cacheGetVal.(domain.Contents), tt.cacheGetErr)

			adapter := NewContentsCacheAdapter(mockCache)

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
		input     *domain.Contents
		setErr    error
		expectErr bool
	}{
		{
			name: "success: valid contents",
			path: "/foo/bar.txt",
			input: &domain.Contents{
				CTag: "etag123",
				Data: []byte("hello"),
			},
			setErr:    nil,
			expectErr: false,
		},
		{
			name: "error: cache.Set fails",
			path: "/foo/fail",
			input: &domain.Contents{
				CTag: "etag",
				Data: []byte("data"),
			},
			setErr:    errors.New("cache write failed"),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[domain.Contents])

			if tt.input != nil {
				mockCache.
					On("Set", mock.Anything, tt.path, *tt.input).
					Return(tt.setErr)
			}

			adapter := NewContentsCacheAdapter(mockCache)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockCache := new(MockCache2[domain.Contents])

			mockCache.
				On("Delete", mock.Anything, tt.path).
				Return(tt.deleteErr)

			adapter := NewContentsCacheAdapter(mockCache)

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
