package util

import (
	"bytes"
	"io"

	applogging "github.com/michaeldcanady/go-onedrive/internal/core/logger/app"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	didomain "github.com/michaeldcanady/go-onedrive/internal/di/domain"
)

// EnsureLogger retrieves or creates the CLI log.
func EnsureLogger(c didomain.Container, name string) (domainlogger.Logger, error) {
	log, err := c.Logger().GetLogger(name)
	if err == applogging.ErrUnknownLogger {
		return c.Logger().CreateLogger(name)
	}
	return log, err
}

// NewReader returns an io.Reader from a byte slice.
func NewReader(b []byte) io.Reader {
	return bytes.NewReader(b)
}

type nopWriteCloser struct {
	io.Writer
}

func (n *nopWriteCloser) Close() error {
	return nil
}

// NewNopWriteCloser returns a WriteCloser with a no-op Close method wrapping the provided Writer.
func NewNopWriteCloser(w io.Writer) io.WriteCloser {
	return &nopWriteCloser{w}
}
