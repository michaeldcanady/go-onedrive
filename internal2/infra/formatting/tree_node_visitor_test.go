package formatting

import (
	"testing"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	"github.com/stretchr/testify/assert"
)

func TestTreeNodeVisitor(t *testing.T) {
	tests := []struct {
		Name   string
		Input  *treeNode
		Result string
	}{
		{
			Name: "Simple",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFile,
					Size: 0,
				},
			},
			Result: "/",
		},
		{
			Name: "1 level",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFolder,
					Size: 0,
				},
				Children: []*treeNode{
					&treeNode{
						Name: "file",
						Item: &fs.Item{
							ID:   "",
							Name: "file",
							Path: "/",
							Type: fs.ItemTypeFolder,
							Size: 0,
						},
					},
				},
			},
			Result: "/\n└──file",
		},
		{
			Name: "2 levels, 1 items",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFolder,
					Size: 0,
				},
				Children: []*treeNode{
					&treeNode{
						Name: "folder/",
						Item: &fs.Item{
							ID:   "",
							Name: "folder",
							Path: "/",
							Type: fs.ItemTypeFolder,
							Size: 0,
						},
						Children: []*treeNode{
							&treeNode{
								Name: "folder/file",
								Item: &fs.Item{
									ID:   "",
									Name: "file",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
						},
					},
				},
			},
			Result: "/\n└──folder\n    └──file",
		},
		{
			Name: "2 levels, 2 items",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFolder,
					Size: 0,
				},
				Children: []*treeNode{
					&treeNode{
						Name: "folder/",
						Item: &fs.Item{
							ID:   "",
							Name: "folder",
							Path: "/",
							Type: fs.ItemTypeFolder,
							Size: 0,
						},
						Children: []*treeNode{
							&treeNode{
								Name: "folder/file1",
								Item: &fs.Item{
									ID:   "",
									Name: "file1",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
							&treeNode{
								Name: "folder/file2",
								Item: &fs.Item{
									ID:   "",
									Name: "file2",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
						},
					},
				},
			},
			Result: "/\n└──folder\n    ├──file1\n    └──file2",
		},
		{
			Name: "2 levels, 2 items",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFolder,
					Size: 0,
				},
				Children: []*treeNode{
					&treeNode{
						Name: "folder/",
						Item: &fs.Item{
							ID:   "",
							Name: "folder",
							Path: "/",
							Type: fs.ItemTypeFolder,
							Size: 0,
						},
						Children: []*treeNode{
							&treeNode{
								Name: "folder/file1",
								Item: &fs.Item{
									ID:   "",
									Name: "file1",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
							&treeNode{
								Name: "folder/file2",
								Item: &fs.Item{
									ID:   "",
									Name: "file2",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
						},
					},
					&treeNode{
						Name: "file",
						Item: &fs.Item{
							ID:   "",
							Name: "file",
							Path: "/",
							Type: fs.ItemTypeFile,
							Size: 0,
						},
					},
				},
			},
			Result: "/\n├──folder\n│   ├──file1\n│   └──file2\n└──file",
		},
		{
			Name: "3 levels, 2 items",
			Input: &treeNode{
				Name: "/",
				Item: &fs.Item{
					ID:   "",
					Name: "/",
					Path: ".",
					Type: fs.ItemTypeFolder,
					Size: 0,
				},
				Children: []*treeNode{
					&treeNode{
						Name: "folder/",
						Item: &fs.Item{
							ID:   "",
							Name: "folder",
							Path: "/",
							Type: fs.ItemTypeFolder,
							Size: 0,
						},
						Children: []*treeNode{
							&treeNode{
								Name: "folder/file1",
								Item: &fs.Item{
									ID:   "",
									Name: "file1",
									Path: "/folder",
									Type: fs.ItemTypeFile,
									Size: 0,
								},
							},
							&treeNode{
								Name: "folder/folder",
								Item: &fs.Item{
									ID:   "",
									Name: "folder",
									Path: "/folder",
									Type: fs.ItemTypeFolder,
									Size: 0,
								},
								Children: []*treeNode{
									&treeNode{
										Name: "file",
										Item: &fs.Item{
											ID:   "",
											Name: "file",
											Path: "/",
											Type: fs.ItemTypeFile,
											Size: 0,
										},
									},
								},
							},
						},
					},
					&treeNode{
						Name: "file",
						Item: &fs.Item{
							ID:   "",
							Name: "file",
							Path: "/",
							Type: fs.ItemTypeFile,
							Size: 0,
						},
					},
				},
			},
			Result: "/\n├──folder\n│   ├──file1\n│   └──folder\n│       └──file\n└──file",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			visitor := &TreeNodeVisitor{}
			visitor.VisitNode(test.Input)
			assert.Equal(t, test.Result, visitor.String())
		})
	}
}
