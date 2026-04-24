package cat

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Execute(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(m *mockFsService)
		wantErr bool
	}{
		{
			name: "successful read",
			setup: func(m *mockFsService) {
				m.On("ReadFile", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(bytes.NewBufferString("test")), nil)
			},
			wantErr: false,
		},
		{
			name: "read failed",
			setup: func(m *mockFsService) {
				m.On("ReadFile", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(bytes.NewBufferString("test")), errors.New("read failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReader := new(mockFsService)
			mockFactory := new(mockURIFactory)
			mockLogger := new(mockLogger)

			handler := NewCommand(mockReader, mockFactory, mockLogger)

			ctx := &CommandContext{
				Ctx: context.Background(),
				URI: &pkgfs.URI{},
				Options: Options{
					Path:   "od:/file",
					Stdout: new(bytes.Buffer),
				},
			}

			mockLogger.On("WithContext", mock.Anything).Return(mockLogger)
			mockLogger.On("With", mock.Anything).Return(mockLogger)
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
			if tt.wantErr {
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
			}

			tt.setup(mockReader)

			err := handler.Execute(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockReader.AssertExpectations(t)
		})
	}
}

type mockURIFactory struct{ mock.Mock }
func (m *mockURIFactory) FromString(s string) (*fsdomain.URI, error) {
	args := m.Called(s)
	return args.Get(0).(*fsdomain.URI), args.Error(1)
}
