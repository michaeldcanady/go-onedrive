package formatting

import domainfs "github.com/michaeldcanady/go-onedrive/internal/fs/shared/domain"

type treeNode struct {
	Name     string
	Item     *domainfs.Item
	Children []*treeNode
}
