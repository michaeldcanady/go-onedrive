package list

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMountLister struct{ mock.Mock }

func (m *MockMountLister) ListMounts(ctx context.Context) ([]mount.MountConfig, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]mount.MountConfig), args.Error(1)
}

type MockFormatCreator struct{ mock.Mock }

func (m *MockFormatCreator) Create(f formatting.Format) (formatting.OutputFormatter, error) {
	args := m.Called(f)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(formatting.OutputFormatter), args.Error(1)
}

type MockOutputFormatter struct{ mock.Mock }

func (m *MockOutputFormatter) Format(w io.Writer, items []any) error {
	return m.Called(w, items).Error(0)
}

func TestHandler_Validate(t *testing.T) {
	mountSvc := new(MockMountLister)
	ff := new(MockFormatCreator)
	cmd := NewCommand(mountSvc, ff, nil)

	ctx := &CommandContext{Options: &Options{Format: "json"}}

	err := cmd.Validate(ctx)
	assert.NoError(t, err)
	assert.Equal(t, formatting.FormatJSON, ctx.Format)
}

func TestHandler_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mountSvc := new(MockMountLister)
		ff := new(MockFormatCreator)
		mockLog := new(MockLogger)
		formatter := new(MockOutputFormatter)
		cmd := NewCommand(mountSvc, ff, mockLog)

		ctx := &CommandContext{
			Ctx:     context.Background(),
			Options: &Options{Stdout: io.Discard},
			Format:  formatting.FormatJSON,
		}

		items := []mount.MountConfig{{Path: "/path"}}
		mountSvc.On("ListMounts", ctx.Ctx).Return(items, nil)
		ff.On("Create", ctx.Format).Return(formatter, nil)
		formatter.On("Format", ctx.Options.Stdout, mock.Anything).Return(nil)
		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount list operation", mock.Anything).Return()
		mockLog.On("Info", "mount list completed successfully", mock.Anything).Return()

		err := cmd.Execute(ctx)
		assert.NoError(t, err)
	})

	t.Run("list failure", func(t *testing.T) {
		mountSvc := new(MockMountLister)
		mockLog := new(MockLogger)
		cmd := NewCommand(mountSvc, nil, mockLog)

		ctx := &CommandContext{Ctx: context.Background()}
		mountSvc.On("ListMounts", ctx.Ctx).Return(nil, errors.New("list failed"))
		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount list operation", mock.Anything).Return()
		mockLog.On("Error", "list failed", mock.Anything).Return()

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
func (m *MockLogger) SetLevel(level logger.Level)             { m.Called(level) }
func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}
func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}
