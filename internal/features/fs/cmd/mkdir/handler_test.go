package mkdir

import (
	"context"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockItemCreator struct{ mock.Mock }
func (m *mockItemCreator) Mkdir(ctx context.Context, uri *fsdomain.URI) error {
	return m.Called(ctx, uri).Error(0)
}

type mockURIFactory struct{ mock.Mock }
func (m *mockURIFactory) FromString(s string) (*fsdomain.URI, error) {
	args := m.Called(s)
	return args.Get(0).(*fsdomain.URI), args.Error(1)
}


func TestHandler_Execute(t *testing.T) {
	mockManager := new(mockItemCreator)
	mockFactory := new(mockURIFactory)
	mockLogger := new(mockLogger)

	handler := NewCommand(mockManager, mockFactory, mockLogger)

	ctx := &CommandContext{
		Ctx: context.Background(),
		Options: Options{
			URI: &fsdomain.URI{},
		},
	}

	mockLogger.On("WithContext", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockManager.On("Mkdir", mock.Anything, mock.Anything).Return(nil)

	err := handler.Execute(ctx)
	assert.NoError(t, err)
	mockManager.AssertExpectations(t)
}
