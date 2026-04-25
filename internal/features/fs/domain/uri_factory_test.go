package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURIFactory_FromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		setup    func(m *mockVFS)
		want     *URI
		wantErr  bool
		errMsg   string
	}{
		{
			name:  "absolute path resolved by VFS",
			input: "/od/test.txt",
			setup: func(m *mockVFS) {
				m.On("Resolve", "/od/test.txt").Return("/od", "/test.txt", nil)
			},
			want: &URI{
				Provider: "/od",
				Path:     "/test.txt",
			},
		},
		{
			name:  "provider prefix resolved by VFS",
			input: "od:/test.txt",
			setup: func(m *mockVFS) {
				m.On("Resolve", "/od").Return("/od", "/", nil)
			},
			want: &URI{
				Provider: "/od",
				Path:     "/test.txt",
			},
		},
		{
			name:  "unknown mount point",
			input: "unknown:/test.txt",
			setup: func(m *mockVFS) {
				m.On("Resolve", "/unknown").Return("", "", assert.AnError)
			},
			wantErr: true,
			errMsg:  "unknown mount point: unknown",
		},
		{
			name:  "default provider fallback",
			input: "local_file.txt",
			want: &URI{
				Provider: DefaultProviderPrefix,
				Path:     "/local_file.txt",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mVFS := new(mockVFS)
			if tt.setup != nil {
				tt.setup(mVFS)
			}

			f := NewURIFactory(mVFS)
			got, err := f.FromString(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			mVFS.AssertExpectations(t)
		})
	}
}

func TestURIFactory_FromLocalPath(t *testing.T) {
	f := NewURIFactory(nil)
	got, err := f.FromLocalPath("test.txt")
	assert.NoError(t, err)
	assert.Equal(t, &URI{Provider: "/local", Path: "/test.txt"}, got)
}

func TestURIFactory_FromMount(t *testing.T) {
	mVFS := new(mockVFS)
	mVFS.On("Resolve", "/od").Return("/od", "/", nil)
	
	f := NewURIFactory(mVFS)
	got, err := f.FromMount("od", "test.txt")
	assert.NoError(t, err)
	assert.Equal(t, &URI{Provider: "/od", Path: "/test.txt"}, got)
	mVFS.AssertExpectations(t)
}
