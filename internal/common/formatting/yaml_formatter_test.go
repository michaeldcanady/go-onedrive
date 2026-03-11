package formatting

import (
	"bytes"
	"testing"

	fs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestYAMLFormatter_Format(t *testing.T) {
	f := &YAMLFormatter{}
	items := []fs.Item{
		{Name: "file1", Size: 100},
		{Name: "folder1", Type: fs.ItemTypeFolder},
	}

	buf := new(bytes.Buffer)
	err := f.Format(buf, items)
	assert.NoError(t, err)

	var result []fs.Item
	err = yaml.Unmarshal(buf.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, len(items), len(result))
	assert.Equal(t, items[0].Name, result[0].Name)
}
