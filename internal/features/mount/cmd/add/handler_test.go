package add

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/identity"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	"github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMountAdder struct{ mock.Mock }

func (m *MockMountAdder) AddMount(ctx context.Context, cfg mount.MountConfig) error {
	return m.Called(ctx, cfg).Error(0)
}

type MockAccountGetter struct{ mock.Mock }

func (m *MockAccountGetter) GetAccount(ctx context.Context, id string) (*identity.Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*identity.Account), args.Error(1)
}

type MockURIFactory struct{ mock.Mock }

func (m *MockURIFactory) FromString(s string) (*fs.URI, error) {
	args := m.Called(s)
	return args.Get(0).(*fs.URI), args.Error(1)
}

func TestHandler_Validate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mountSvc := new(MockMountAdder)
		accountSvc := new(MockAccountGetter)
		uriFactory := new(MockURIFactory)
		cmd := NewCommand(mountSvc, accountSvc, uriFactory, nil)

		uri := &fs.URI{}
		uriFactory.On("FromString", "/path").Return(uri, nil)
		accountSvc.On("GetAccount", mock.Anything, "id").Return(&identity.Account{ID: "id"}, nil)

		ctx := &CommandContext{
			Ctx: context.Background(),
			Options: &Options{
				Path:       "/path",
				Type:       "onedrive",
				IdentityID: "id",
			},
		}

		err := cmd.Validate(ctx)
		assert.NoError(t, err)
		assert.Equal(t, uri, ctx.Uri)
		assert.Equal(t, "onedrive", ctx.Type)
	})

	t.Run("invalid path", func(t *testing.T) {
		uriFactory := new(MockURIFactory)
		cmd := NewCommand(nil, nil, uriFactory, nil)
		uriFactory.On("FromString", "bad").Return((*fs.URI)(nil), errors.New("parse error"))

		ctx := &CommandContext{Options: &Options{Path: "bad"}}
		err := cmd.Validate(ctx)
		assert.Error(t, err)
	})
}

func TestHandler_Execute(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mountSvc := new(MockMountAdder)
		mockLog := new(MockLogger) // Need a MockLogger here too
		cmd := NewCommand(mountSvc, nil, nil, mockLog)

		ctx := &CommandContext{
			Ctx:          context.Background(),
			Uri:          &fs.URI{},
			Type:         "onedrive",
			Identity:     &identity.Account{ID: "id"},
			MountOptions: map[string]string{"key": "value"},
		}

		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount add operation", mock.Anything).Return()
		mockLog.On("Info", "mount add completed successfully", mock.Anything).Return()
		mountSvc.On("AddMount", ctx.Ctx, mount.MountConfig{
			Path:       ctx.Uri.String(),
			Type:       ctx.Type,
			IdentityID: ctx.Identity.ID,
			Options:    ctx.MountOptions,
		}).Return(nil)

		err := cmd.Execute(ctx)
		assert.NoError(t, err)
		mountSvc.AssertExpectations(t)
	})

	t.Run("add mount failure", func(t *testing.T) {
		mountSvc := new(MockMountAdder)
		mockLog := new(MockLogger)
		cmd := NewCommand(mountSvc, nil, nil, mockLog)

		ctx := &CommandContext{
			Ctx:      context.Background(),
			Uri:      &fs.URI{},
			Type:     "onedrive",
			Identity: &identity.Account{ID: "id"},
		}

		mockLog.On("WithContext", ctx.Ctx).Return(mockLog)
		mockLog.On("Info", "starting mount add operation", mock.Anything).Return()
		mockLog.On("Error", "mount add failed", mock.Anything).Return()
		mountSvc.On("AddMount", mock.Anything, mock.Anything).Return(errors.New("add failed"))

		err := cmd.Execute(ctx)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to add mount")
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
