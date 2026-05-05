package remove

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMountRemover struct{ mock.Mock }

func (m *MockMountRemover) RemoveMount(ctx context.Context, path string) error {
	return m.Called(ctx, path).Error(0)
}

type MockURIFactory struct{ mock.Mock }

func (m *MockURIFactory) FromString(s string) (*fs.URI, error) {
	args := m.Called(s)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*fs.URI), args.Error(1)
}

func TestHandler_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mountSvc := new(MockMountRemover)
		uriFactory := new(MockURIFactory)
		cmd := NewCommand(mountSvc, uriFactory, nil)

		uri := &fs.URI{}
		uriFactory.On("FromString", "/path").Return(uri, nil)

		ctx := &CommandContext{
			Ctx: context.Background(),
			Options: &Options{
				Path: "/path",
			},
		}

		err := cmd.Validate(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uri, ctx.Uri)
	})

	t.Run("invalid path", func(t *testing.T) {
		uriFactory := new(MockURIFactory)
		cmd := NewCommand(nil, uriFactory, nil)
		uriFactory.On("FromString", "bad").Return((*fs.URI)(nil), errors.New("parse error"))

		ctx := &CommandContext{Options: &Options{Path: "bad"}}
		err := cmd.Validate(ctx)
		assert.Error(t, err)
	})
}

func TestHandler_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mountSvc := new(MockMountRemover)
		mockLog := new(MockLogger)
		cmd := NewCommand(mountSvc, nil, mockLog)

		ctx := &CommandContext{
			Ctx: context.Background(),
			Uri: &fs.URI{},
		}

		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount remove operation", mock.Anything).Return()
		mockLog.On("Info", "mount remove completed successfully", mock.Anything).Return()
		mountSvc.On("RemoveMount", ctx.Ctx, mock.Anything).Return(nil)

		err := cmd.Execute(ctx)
		assert.NoError(t, err)
	})

	t.Run("remove failure", func(t *testing.T) {
		mountSvc := new(MockMountRemover)
		mockLog := new(MockLogger)
		cmd := NewCommand(mountSvc, nil, mockLog)

		ctx := &CommandContext{
			Ctx: context.Background(),
			Uri: &fs.URI{},
		}

		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount remove operation", mock.Anything).Return()
		mountSvc.On("RemoveMount", mock.Anything, mock.Anything).Return(errors.New("remove failed"))

		err := cmd.Execute(ctx)
		assert.Error(t, err)
	})
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *MockLogger) Info(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Warn(msg string, fields ...logger.Field)  { m.Called(msg, fields) }
func (m *MockLogger) Error(msg string, fields ...logger.Field) { m.Called(msg, fields) }
func (m *MockLogger) SetLevel(level logger.Level)              { m.Called(level) }
func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}
func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}
