package upload

import (
	"context"
	"errors"
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
			name: "successful upload",
			setup: func(m *mockFsService) {
				m.On("Copy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "upload failed",
			setup: func(m *mockFsService) {
				m.On("Copy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("upload failed"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockManager := new(mockFsService)
			mockFactory := fsdomain.NewURIFactory(new(mockVFS))
			mockLogger := new(mockLogger)

			handler := NewCommand(mockManager, mockFactory, mockLogger)

			ctx := &CommandContext{
				Ctx: context.Background(),
				Options: Options{
					Source:         "/local/src",
					SourceURI:      &pkgfs.URI{},
					Destination:    "/od/dst",
					DestinationURI: &pkgfs.URI{},
				},
			}

			mockLogger.On("WithContext", mock.Anything).Return(mockLogger)
			mockLogger.On("With", mock.Anything).Return(mockLogger)
			mockLogger.On("Info", mock.Anything, mock.Anything).Return()
			mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
			if tt.wantErr {
				mockLogger.On("Error", mock.Anything, mock.Anything).Return()
			}

			tt.setup(mockManager)

			err := handler.Execute(ctx)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			mockManager.AssertExpectations(t)
		})
	}
}
