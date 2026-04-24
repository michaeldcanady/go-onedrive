package download

import (
	"context"
	"testing"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Execute(t *testing.T) {
	mockManager := new(mockFsService)
	mockFactory := fs.NewURIFactory(new(mockVFS))
	mockLogger := new(mockLogger)

	handler := NewCommand(mockManager, mockFactory, mockLogger)

	ctx := &CommandContext{
		Ctx:        context.Background(),
		SourceURI:      &pkgfs.URI{},
		DestinationURI: &pkgfs.URI{},
		Options: Options{
			Source:      "od:/src",
			Destination: "/local/dst",
		},
	}

	mockLogger.On("WithContext", mock.Anything).Return(mockLogger)
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockManager.On("Copy", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := handler.Execute(ctx)
	assert.NoError(t, err)
	mockManager.AssertExpectations(t)
}
