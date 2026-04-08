package drive

import (
	"context"
	"errors"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"github.com/michaeldcanady/go-onedrive/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGateway is a mock implementation of Gateway.
type MockGateway struct {
	mock.Mock
}

func (m *MockGateway) ListDrives(ctx context.Context) ([]Drive, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Drive), args.Error(1)
}

func (m *MockGateway) GetPersonalDrive(ctx context.Context) (Drive, error) {
	args := m.Called(ctx)
	return args.Get(0).(Drive), args.Error(1)
}

// MockStateService is a mock implementation of state.Service.
type MockStateService struct {
	mock.Mock
}

func (m *MockStateService) Get(key state.Key) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockStateService) Set(key state.Key, value string, scope state.Scope) error {
	args := m.Called(key, value, scope)
	return args.Error(0)
}

func (m *MockStateService) Clear(key state.Key) error {
	args := m.Called(key)
	return args.Error(0)
}

// MockLogger is a mock implementation of logger.Logger.
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Info(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Warn(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Error(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) Debug(msg string, kv ...logger.Field) {
	m.Called(msg, kv)
}

func (m *MockLogger) SetLevel(level logger.Level) {
	m.Called(level)
}

func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
	return m
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	return m
}

func TestDefaultService(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		setup     func(g *MockGateway, s *MockStateService)
		operation func(t *testing.T, s *DefaultService)
	}{
		{
			name: "ListDrives success",
			setup: func(g *MockGateway, s *MockStateService) {
				g.On("ListDrives", ctx).Return([]Drive{{ID: "d1"}}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				drives, err := s.ListDrives(ctx)
				assert.NoError(t, err)
				assert.Len(t, drives, 1)
				assert.Equal(t, "d1", drives[0].ID)
			},
		},
		{
			name: "ResolveDrive by ID",
			setup: func(g *MockGateway, s *MockStateService) {
				g.On("ListDrives", ctx).Return([]Drive{{ID: "d1", Name: "N1"}}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				d, err := s.ResolveDrive(ctx, "d1")
				assert.NoError(t, err)
				assert.Equal(t, "d1", d.ID)
			},
		},
		{
			name: "ResolveDrive by Name",
			setup: func(g *MockGateway, s *MockStateService) {
				g.On("ListDrives", ctx).Return([]Drive{{ID: "d1", Name: "N1"}}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				d, err := s.ResolveDrive(ctx, "N1")
				assert.NoError(t, err)
				assert.Equal(t, "d1", d.ID)
			},
		},
		{
			name: "ResolvePersonalDrive success",
			setup: func(g *MockGateway, s *MockStateService) {
				g.On("GetPersonalDrive", ctx).Return(Drive{ID: "p1"}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				d, err := s.ResolvePersonalDrive(ctx)
				assert.NoError(t, err)
				assert.Equal(t, "p1", d.ID)
			},
		},
		{
			name: "GetActive from state",
			setup: func(g *MockGateway, s *MockStateService) {
				s.On("Get", state.KeyDrive).Return("d1", nil).Once()
				g.On("ListDrives", ctx).Return([]Drive{{ID: "d1"}}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				d, err := s.GetActive(ctx)
				assert.NoError(t, err)
				assert.Equal(t, "d1", d.ID)
			},
		},
		{
			name: "GetActive fallback to personal",
			setup: func(g *MockGateway, s *MockStateService) {
				s.On("Get", state.KeyDrive).Return("", nil).Once()
				g.On("GetPersonalDrive", ctx).Return(Drive{ID: "p1"}, nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				d, err := s.GetActive(ctx)
				assert.NoError(t, err)
				assert.Equal(t, "p1", d.ID)
			},
		},
		{
			name: "SetActive updates state",
			setup: func(g *MockGateway, s *MockStateService) {
				s.On("Set", state.KeyDrive, "d1", state.ScopeGlobal).Return(nil).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				err := s.SetActive(ctx, "d1", state.ScopeGlobal)
				assert.NoError(t, err)
			},
		},
		{
			name: "Gateway error propagation",
			setup: func(g *MockGateway, s *MockStateService) {
				g.On("ListDrives", ctx).Return(nil, errors.New("fail")).Once()
			},
			operation: func(t *testing.T, s *DefaultService) {
				_, err := s.ListDrives(ctx)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "fail")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockGateway := new(MockGateway)
			mockState := new(MockStateService)
			mockLog := new(MockLogger)
			mockLog.On("Debug", mock.Anything, mock.Anything).Maybe()
			mockLog.On("Info", mock.Anything, mock.Anything).Maybe()
			mockLog.On("Warn", mock.Anything, mock.Anything).Maybe()
			mockLog.On("Error", mock.Anything, mock.Anything).Maybe()

			if tt.setup != nil {
				tt.setup(mockGateway, mockState)
			}

			service := NewDefaultService(mockGateway, mockState, mockLog)
			tt.operation(t, service)

			mockGateway.AssertExpectations(t)
			mockState.AssertExpectations(t)
		})
	}
}
