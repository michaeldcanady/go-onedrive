package cat

import (
	"bytes"
	"context"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandler_Execute(t *testing.T) {
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
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockReader.On("ReadFile", mock.Anything, mock.Anything, mock.Anything).Return(io.NopCloser(bytes.NewBufferString("content")), nil)

	err := handler.Execute(ctx)
	assert.NoError(t, err)
	mockReader.AssertExpectations(t)
}

type mockURIFactory struct{ mock.Mock }
func (m *mockURIFactory) FromString(s string) (*fsdomain.URI, error) {
	args := m.Called(s)
	return args.Get(0).(*fsdomain.URI), args.Error(1)
}
