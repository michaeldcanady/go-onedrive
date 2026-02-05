package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"

type treeNode struct {
	Name     string
	Item     *domainfs.Item
	Children []*treeNode
}
