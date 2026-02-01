package formatting

import (
	"fmt"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type TreeNodeVisitor struct {
	builder   strings.Builder
	prefix    string
	connector string
	level     int
}

func (v *TreeNodeVisitor) VisitNode(node *treeNode) error {
	switch node.Item.Type {
	case fs.ItemTypeFolder:
		v.VisitFolder(node)
	case fs.ItemTypeFile:
		v.VisitFile(node)
	default:
		return fmt.Errorf("unknown item type: %d", node.Item.Type)
	}
	return nil
}

func (v *TreeNodeVisitor) VisitFolder(node *treeNode) error {
	childrenCount := len(node.Children)

	v.builder.WriteString(v.prefix + v.connector + node.Item.Name)

	for i, child := range node.Children {
		v.builder.WriteString("\n")
		isLast := i == childrenCount-1

		// Set connector for this child
		prevConnector := v.connector
		if isLast {
			v.connector = "└──"
		} else {
			v.connector = "├──"
		}

		// Save prefix
		prevPrefix := v.prefix

		// Extend prefix based on *parent* connector
		switch prevConnector {
		case "├──":
			// Parent has more siblings → vertical bar continues
			v.prefix = prevPrefix + "│   "
		case "└──":
			// Parent was last → no vertical bar
			v.prefix = prevPrefix + "    "
		}

		// Recurse
		if err := v.VisitNode(child); err != nil {
			return err
		}

		// Restore
		v.prefix = prevPrefix
		v.connector = prevConnector
	}
	return nil
}

func (v *TreeNodeVisitor) VisitFile(node *treeNode) {
	v.builder.WriteString(v.prefix + v.connector + node.Item.Name)
}

func (v *TreeNodeVisitor) String() string {
	return v.builder.String()
}
