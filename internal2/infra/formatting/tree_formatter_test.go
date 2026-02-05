package formatting

import (
	"testing"

	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/stretchr/testify/assert"
)

func TestTreeFormatter_buildTree(t *testing.T) {
	tests := []struct {
		Name     string
		Input    []domainfs.Item
		Expected *treeNode
	}{
		{
			Name:     "Empty",
			Input:    nil,
			Expected: nil,
		},
		{
			Name: "1 element",
			Input: []domainfs.Item{
				{
					Name: "file",
					Path: "/",
				},
			},
			Expected: &treeNode{
				Name: "file",
				Item: &domainfs.Item{
					Name: "file",
					Path: "/",
				},
				Children: []*treeNode{},
			},
		},
		{
			Name: "depth 1",
			Input: []domainfs.Item{
				{
					Name: "file1",
					Path: "/",
					Type: domainfs.ItemTypeFile,
				},
				{
					Name: "file2",
					Path: "/",
					Type: domainfs.ItemTypeFile,
				},
			},
			Expected: &treeNode{
				Name: "/",
				Item: &domainfs.Item{
					Name: "/",
					Path: "",
					Type: domainfs.ItemTypeFolder,
				},
				Children: []*treeNode{
					{
						Name: "file1",
						Item: &domainfs.Item{
							Name: "file1",
							Path: "/",
						},
						Children: []*treeNode{},
					},
					{
						Name: "file2",
						Item: &domainfs.Item{
							Name: "file2",
							Path: "/",
						},
						Children: []*treeNode{},
					},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := NewTreeFormatter().buildTree(test.Input)

			assert.Equal(t, test.Expected, output)
		})
	}
}
