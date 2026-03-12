package formatting

import (
	"github.com/michaeldcanady/go-onedrive/internal2/core/fs/shared"
)

// treeNode represents a single element in a hierarchical filesystem tree.
type treeNode struct {
	// Name is the display label for the node.
	Name string
	// Item is the underlying filesystem object (may be nil for synthetic root).
	Item *shared.Item
	// Children contains the nested nodes within this element.
	Children []*treeNode
}
