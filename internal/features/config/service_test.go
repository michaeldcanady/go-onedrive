package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRepository struct {
	mock.Mock
}

func (m *mockRepository) Load() (map[string]any, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]any), args.Error(1)
}

func (m *mockRepository) Save(config map[string]any) error {
	args := m.Called(config)
	return args.Error(0)
}

func TestConfigService_Get(t *testing.T) {
	tests := []struct {
		name    string
		config  map[string]any
		key     string
		want    any
		wantErr bool
		errMsg  string
	}{
		{
			name: "get string value",
			config: map[string]any{
				"core": map[string]any{
					"plugins_dir": "/tmp/plugins",
				},
			},
			key:     "core.plugins_dir",
			want:    "/tmp/plugins",
			wantErr: false,
		},
		{
			name: "get map value (raw)",
			config: map[string]any{
				"identity": map[string]any{
					"google": map[string]any{
						"client_id":     "google-id",
						"client_secret": "google-secret",
					},
				},
			},
			key: "identity.google",
			want: map[string]any{
				"client_id":     "google-id",
				"client_secret": "google-secret",
			},
			wantErr: false,
		},
		{
			name:    "get default value",
			config:  map[string]any{},
			key:     "identity.azure.client_id",
			want:    "6b1e6ec0-ad93-4175-a0e0-84c02e13f206",
			wantErr: false,
		},
		{
			name:    "key not found",
			config:  map[string]any{},
			key:     "non.existent.key",
			wantErr: true,
			errMsg:  "key not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(mockRepository)
			repo.On("Load").Return(tt.config, nil)

			s := NewConfigService(repo, nil)
			got, err := s.Get(tt.key)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
				// Critical check: Ensure it's not a string if it's a map
				if _, ok := tt.want.(map[string]any); ok {
					_, isString := got.(string)
					assert.False(t, isString, "Result should be a map, not a string representation")
				}
			}
		})
	}
}
