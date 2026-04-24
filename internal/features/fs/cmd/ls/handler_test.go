package ls

import (
	"bytes"
	"context"
	"io"
	"testing"

	fsdomain "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockItemLister struct{ mock.Mock }
func (m *mockItemLister) List(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ListOptions) ([]pkgfs.Item, error) {
	args := m.Called(ctx, uri, opts)
	return args.Get(0).([]pkgfs.Item), args.Error(1)
}

type mockFormatCreator struct{ mock.Mock }
func (m *mockFormatCreator) Create(f formatting.Format) (formatting.OutputFormatter, error) {
	args := m.Called(f)
	return args.Get(0).(formatting.OutputFormatter), args.Error(1)
}

type mockFormatter struct{ mock.Mock }
func (m *mockFormatter) Format(w io.Writer, items []any) error { return m.Called(w, items).Error(0) }

type mockURIFactory struct{ mock.Mock }
func (m *mockURIFactory) FromString(s string) (*fsdomain.URI, error) {
	args := m.Called(s)
	return args.Get(0).(*fsdomain.URI), args.Error(1)
}


func TestHandler_Execute(t *testing.T) {
	mockManager := new(mockItemLister)
	mockFactory := new(mockURIFactory)
	mockFormatCreator := new(mockFormatCreator)
	mockLogger := new(mockLogger)
	mockFormatter := new(mockFormatter)

	handler := NewCommand(mockManager, mockFactory, mockFormatCreator, mockLogger)

	ctx := &CommandContext{
		Ctx: context.Background(),
		Options: Options{
			URI:    &fsdomain.URI{},
			Format: formatting.FormatShort,
			Stdout: new(bytes.Buffer),
		},
	}

	mockLogger.On("WithContext", mock.Anything).Return(mockLogger)
	mockLogger.On("With", mock.Anything).Return(mockLogger)
	mockLogger.On("Info", mock.Anything, mock.Anything).Return()
	mockLogger.On("Debug", mock.Anything, mock.Anything).Return()
	mockManager.On("List", mock.Anything, mock.Anything, mock.Anything).Return([]pkgfs.Item{}, nil)
	mockFormatCreator.On("Create", formatting.FormatShort).Return(mockFormatter, nil)
	mockFormatter.On("Format", mock.Anything, mock.Anything).Return(nil)

	err := handler.Execute(ctx)
	assert.NoError(t, err)
	mockManager.AssertExpectations(t)
}
