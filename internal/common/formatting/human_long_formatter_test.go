package formatting

import (
	"bytes"
	"strings"
	"testing"
	"time"

	domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/domain"
	"github.com/stretchr/testify/assert"
)

func TestHumanLongFormatter_Format(t *testing.T) {
	formatter := &HumanLongFormatter{}
	buf := new(bytes.Buffer)

	now := time.Now()
	items := []domainfs.Item{
		{
			Name:     "file.txt",
			Path:     "/path/to",
			Type:     domainfs.ItemTypeFile,
			Size:     1024,
			Modified: now,
		},
		{
			Name:     "folder",
			Path:     "/path/to",
			Type:     domainfs.ItemTypeFolder,
			Modified: now,
		},
	}

	err := formatter.Format(buf, items)
	assert.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "Modified")
	assert.Contains(t, output, "Size")
	assert.Contains(t, output, "Name")
	assert.Contains(t, output, "file.txt")
	assert.Contains(t, output, "folder")
	assert.Contains(t, output, "1.00 KiB")
	assert.Contains(t, output, "-")

	lines := strings.Split(strings.TrimSpace(output), "\n")
	assert.Equal(t, 4, len(lines)) // Header + Separator + 2 Items
}
