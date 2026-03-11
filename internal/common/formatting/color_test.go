package formatting

import (
	"bytes"
	"testing"

	fs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/stretchr/testify/assert"
)

type mockFd struct {
	fd uintptr
}

func (m mockFd) Fd() uintptr {
	return m.fd
}

func TestIsTerminal(t *testing.T) {
	assert.False(t, IsTerminal(nil))
	assert.False(t, IsTerminal(new(bytes.Buffer)))

	// A mockFd will implement the interface but term.IsTerminal will likely return false for a random fd
	assert.False(t, IsTerminal(mockFd{fd: 999}))
}

func TestColorize(t *testing.T) {
	buf := new(bytes.Buffer)
	text := "hello"

	// Since buf is not a terminal, Colorize should return plain text
	result := Colorize(buf, Red, text)
	assert.Equal(t, text, result)
}

func TestColorizeItem(t *testing.T) {
	buf := new(bytes.Buffer)

	file := fs.Item{Name: "file.txt", Type: fs.ItemTypeFile}
	folder := fs.Item{Name: "folder", Type: fs.ItemTypeFolder}

	// No colors because buf is not a terminal
	assert.Equal(t, "file.txt", ColorizeItem(buf, file))
	assert.Equal(t, "folder", ColorizeItem(buf, folder))
}
