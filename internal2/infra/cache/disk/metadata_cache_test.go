package disk

import (
	"context"
	"errors"
	"testing"

	mocks "github.com/michaeldcanady/go-onedrive/internal2/infra/cache/internal/mock"
	"github.com/stretchr/testify/require"
)

func TestMetadataCache(t *testing.T) {
	type Meta struct {
		Version int
		Note    string
	}

	tests := []struct {
		name       string
		key        string
		meta       Meta
		serialized []byte
		expectErr  bool
	}{
		{
			name:       "set and get metadata",
			key:        "user1",
			meta:       Meta{Version: 1, Note: "hello"},
			serialized: []byte(`{"Version":1,"Note":"hello"}`),
		},
		{
			name:      "serializer error on set",
			key:       "bad",
			meta:      Meta{Version: -1},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Mocks
			mockCache := new(mocks.MockCache[string, string])
			mockSerializer := new(mocks.SerializerMock[Meta])

			// Key serializer mock
			mockKeySerializer := new(mocks.SerializerMock[string])
			mockKeySerializer.
				On("Serialize", tt.key).
				Return([]byte(tt.key), nil)

			mockCache.
				On("KeySerializer").
				Return(mockKeySerializer)

			// Metadata serializer expectations
			if tt.expectErr {
				mockSerializer.
					On("Serialize", tt.meta).
					Return([]byte(nil), errors.New("serialize error"))
			} else {
				mockSerializer.
					On("Serialize", tt.meta).
					Return(tt.serialized, nil)

				mockSerializer.
					On("Deserialize", tt.serialized).
					Return(tt.meta, nil)
			}

			// Create metadata cache
			mc := &MetadataCache[string, string, Meta]{
				Cache:              mockCache,
				metadataIndex:      make(map[string]int64),
				metadataSerializer: mockSerializer,
				metadataPath:       t.TempDir() + "test_metadata.db",
			}

			// Set metadata
			err := mc.SetMetadata(context.Background(), tt.key, tt.meta)
			if tt.expectErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Get metadata
			got, err := mc.GetMetadata(context.Background(), tt.key)
			require.NoError(t, err)
			require.Equal(t, tt.meta, got)

			mockCache.AssertExpectations(t)
			mockSerializer.AssertExpectations(t)
			mockKeySerializer.AssertExpectations(t)
		})
	}
}
