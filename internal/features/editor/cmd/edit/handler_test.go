package edit

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	"github.com/michaeldcanady/go-onedrive/internal/features/vfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockVFS struct {
	mock.Mock
}

func (m *mockVFS) List(ctx context.Context, path string) ([]*vfs.Node, error) {
	args := m.Called(ctx, path)
	return args.Get(0).([]*vfs.Node), args.Error(1)
}

func (m *mockVFS) Stat(ctx context.Context, path string) (*vfs.Node, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*vfs.Node), args.Error(1)
}

func (m *mockVFS) Mkdir(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *mockVFS) Remove(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	return args.Error(0)
}

func (m *mockVFS) Move(ctx context.Context, src, dst string) error {
	args := m.Called(ctx, src, dst)
	return args.Error(0)
}

func (m *mockVFS) Copy(ctx context.Context, src, dst string) error {
	args := m.Called(ctx, src, dst)
	return args.Error(0)
}

func (m *mockVFS) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	args := m.Called(ctx, path)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockVFS) Write(ctx context.Context, path string, reader io.Reader, options ...vfs.WriteOption) error {
	opts := make(map[string]string)
	for _, opt := range options {
		opt(opts)
	}
	args := m.Called(ctx, path, reader, opts)
	return args.Error(0)
}

type mockEditor struct {
	mock.Mock
}

func (m *mockEditor) Open(ctx context.Context, path string) error {
	args := m.Called(ctx, path)
	// Simulate modification by touching the file
	if args.Bool(1) {
		_ = os.Chtimes(path, time.Now(), time.Now().Add(time.Hour))
	}
	return args.Error(0)
}

type mockLogger struct {
	mock.Mock
}

func (m *mockLogger) Debug(msg string, fields ...any) {}
func (m *mockLogger) Info(msg string, fields ...any)  {}
func (m *mockLogger) Warn(msg string, fields ...any)  {}
func (m *mockLogger) Error(msg string, fields ...any) {}
func (m *mockLogger) Fatal(msg string, fields ...any) {}

func (m *mockLogger) With(fields ...any) logger.Service {
	return m
}

func (m *mockLogger) Sync() error {
	return nil
}

func (m *mockLogger) SetLevel(level string) error {
	return nil
}

func (m *mockLogger) GetLevel() string {
	return "info"
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name          string
		force         bool
		initialETag   string
		writeError    error
		expectIfMatch bool
		expectError   bool
	}{
		{
			name:          "successful upload with if-match",
			force:         false,
			initialETag:   "etag123",
			writeError:    nil,
			expectIfMatch: true,
			expectError:   false,
		},
		{
			name:          "failed upload due to etag mismatch",
			force:         false,
			initialETag:   "etag123",
			writeError:    fmt.Errorf("precondition failed"),
			expectIfMatch: true,
			expectError:   true,
		},
		{
			name:          "force upload ignores etag mismatch",
			force:         true,
			initialETag:   "etag123",
			writeError:    nil,
			expectIfMatch: false,
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mvfs := new(mockVFS)
			meditor := new(mockEditor)
			mlog := new(mockLogger)

			ctx := context.Background()
			path := "/test/file.txt"
			node := &vfs.Node{Path: path, ETag: tt.initialETag}

			mvfs.On("Stat", ctx, path).Return(node, nil)
			mvfs.On("Read", ctx, path).Return(io.NopCloser(bytes.NewReader([]byte("content"))), nil)
			meditor.On("Open", ctx, mock.Anything).Return(nil, true) // true means simulate modification

			expectedOpts := make(map[string]string)
			if tt.expectIfMatch {
				expectedOpts["if_match"] = tt.initialETag
			}

			mvfs.On("Write", ctx, path, mock.Anything, expectedOpts).Return(tt.writeError)

			cmd := NewCommand(mvfs, nil, meditor, mlog, mlog)
			opts := Options{
				Path:  path,
				Force: tt.force,
			}

			err := cmd.Execute(&CommandContext{
				Ctx:     ctx,
				Options: opts,
			})

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mvfs.AssertExpectations(t)
			meditor.AssertExpectations(t)
		})
	}
}
