package middleware

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLogger mocks the Logger interface
type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) Debug(msg string, fields ...logger.Field) {
	_ = m.Called(msg, fields)
}

func (m *MockLogger) Error(msg string, fields ...logger.Field) {
	_ = m.Called(msg, fields)
}

func (m *MockLogger) Info(msg string, fields ...logger.Field) {
	_ = m.Called(msg, fields)
}

func (m *MockLogger) Warn(msg string, fields ...logger.Field) {
	_ = m.Called(msg, fields)
}

func (m *MockLogger) SetLevel(level logger.Level) {
	m.Called(level)
}

func (m *MockLogger) With(fields ...logger.Field) logger.Logger {
	args := m.Called(fields)
	return args.Get(0).(logger.Logger)
}

func (m *MockLogger) WithContext(ctx context.Context) logger.Logger {
	args := m.Called(ctx)
	return args.Get(0).(logger.Logger)
}

// MockPipeline mocks the nethttp.Pipeline interface
type MockPipeline struct {
	mock.Mock
}

func (m *MockPipeline) Next(req *http.Request, middlewareIndex int) (*http.Response, error) {
	args := m.Called(req, middlewareIndex)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestKiotaLoggingMiddleware_Intercept(t *testing.T) {
	t.Run("successful request and response logging", func(t *testing.T) {
		mockLog := new(MockLogger)
		mockPipeline := new(MockPipeline)
		middleware := NewKiotaLoggingMiddleware(mockLog)

		req, _ := http.NewRequest("GET", "https://example.com", nil)
		resp := &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(""))}

		mockLog.On("WithContext", req.Context()).Return(mockLog)
		mockLog.On("Debug", "outbound request", mock.Anything).Return()
		mockPipeline.On("Next", req, 0).Return(resp, nil)
		mockLog.On("Debug", "inbound response", mock.Anything).Return()

		result, err := middleware.Intercept(mockPipeline, 0, req)

		assert.NoError(t, err)
		assert.Equal(t, resp, result)
		mockLog.AssertExpectations(t)
		mockPipeline.AssertExpectations(t)
	})

	t.Run("failed request logging", func(t *testing.T) {
		mockLog := new(MockLogger)
		mockPipeline := new(MockPipeline)
		middleware := NewKiotaLoggingMiddleware(mockLog)

		req, _ := http.NewRequest("POST", "https://example.com", nil)
		err := errors.New("network failure")

		mockLog.On("WithContext", req.Context()).Return(mockLog)
		mockLog.On("Debug", "outbound request", mock.Anything).Return()
		mockPipeline.On("Next", req, 0).Return(nil, err)
		mockLog.On("Error", "request failed", mock.Anything).Return()

		result, errResult := middleware.Intercept(mockPipeline, 0, req)

		assert.Error(t, errResult)
		assert.Nil(t, result)
		mockLog.AssertExpectations(t)
		mockPipeline.AssertExpectations(t)
	})
}
