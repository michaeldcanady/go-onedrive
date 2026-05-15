package vfs

import (
	"context"
	"io"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

// Middleware is a function that wraps a [VFS] implementation.
type Middleware func(VFS) VFS

// LoggingMiddleware returns a [Middleware] that logs VFS operations.
func LoggingMiddleware(l logger.Service) Middleware {
	return func(next VFS) VFS {
		return &loggingMiddleware{
			next:   next,
			logger: l,
		}
	}
}

type loggingMiddleware struct {
	next   VFS
	logger logger.Service
}

func (m *loggingMiddleware) List(ctx context.Context, path string) ([]*Node, error) {
	start := time.Now()
	nodes, err := m.next.List(ctx, path)
	m.log(ctx, "List", path, start, err)
	return nodes, err
}

func (m *loggingMiddleware) Stat(ctx context.Context, path string) (*Node, error) {
	start := time.Now()
	node, err := m.next.Stat(ctx, path)
	m.log(ctx, "Stat", path, start, err)
	return node, err
}

func (m *loggingMiddleware) Mkdir(ctx context.Context, path string) error {
	start := time.Now()
	err := m.next.Mkdir(ctx, path)
	m.log(ctx, "Mkdir", path, start, err)
	return err
}

func (m *loggingMiddleware) Remove(ctx context.Context, path string) error {
	start := time.Now()
	err := m.next.Remove(ctx, path)
	m.log(ctx, "Remove", path, start, err)
	return err
}

func (m *loggingMiddleware) Move(ctx context.Context, src, dst string) error {
	start := time.Now()
	err := m.next.Move(ctx, src, dst)
	m.logComplex(ctx, "Move", src, dst, start, err)
	return err
}

func (m *loggingMiddleware) Copy(ctx context.Context, src, dst string) error {
	start := time.Now()
	err := m.next.Copy(ctx, src, dst)
	m.logComplex(ctx, "Copy", src, dst, start, err)
	return err
}

func (m *loggingMiddleware) Read(ctx context.Context, path string) (io.ReadCloser, error) {
	start := time.Now()
	reader, err := m.next.Read(ctx, path)
	m.log(ctx, "Read", path, start, err)
	return reader, err
}

func (m *loggingMiddleware) Write(ctx context.Context, path string, reader io.Reader, options ...WriteOption) error {
	start := time.Now()
	err := m.next.Write(ctx, path, reader, options...)
	m.log(ctx, "Write", path, start, err)
	return err
}

func (m *loggingMiddleware) log(ctx context.Context, op, path string, start time.Time, err error) {
	l := logger.WithContext(m.logger, ctx)
	duration := time.Since(start)

	if err != nil {
		l.Error("VFS operation failed", "op", op, "path", path, "duration", duration, "error", err)
	} else {
		l.Debug("VFS operation succeeded", "op", op, "path", path, "duration", duration)
	}
}

func (m *loggingMiddleware) logComplex(ctx context.Context, op, src, dst string, start time.Time, err error) {
	l := logger.WithContext(m.logger, ctx)
	duration := time.Since(start)

	if err != nil {
		l.Error("VFS operation failed", "op", op, "src", src, "dst", dst, "duration", duration, "error", err)
	} else {
		l.Debug("VFS operation succeeded", "op", op, "src", src, "dst", dst, "duration", duration)
	}
}
