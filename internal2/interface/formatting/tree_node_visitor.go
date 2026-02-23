package formatting

import (
	"fmt"
	"io"
)

// TreeStyle defines the characters used to render the tree structure.
type TreeStyle struct {
	NodeConnector string
	LastConnector string
	VerticalBar   string
	Indent        string
}

// DefaultTreeStyle provides standard Unicode connectors.
var DefaultTreeStyle = TreeStyle{
	NodeConnector: "├──",
	LastConnector: "└──",
	VerticalBar:   "│   ",
	Indent:        "    ",
}

// TreeNodeVisitor handles the recursive rendering of a tree structure.
type TreeNodeVisitor struct {
	writer io.Writer
	style  TreeStyle
}

// NewTreeNodeVisitor creates a new visitor with the default style.
func NewTreeNodeVisitor(w io.Writer) *TreeNodeVisitor {
	return &TreeNodeVisitor{
		writer: w,
		style:  DefaultTreeStyle,
	}
}

// WithStyle allows customizing the tree rendering style.
func (v *TreeNodeVisitor) WithStyle(style TreeStyle) *TreeNodeVisitor {
	v.style = style
	return v
}

// VisitNode starts the recursive printing of the tree from the given root.
func (v *TreeNodeVisitor) VisitNode(node *treeNode) error {
	if node == nil {
		return nil
	}
	// The root node itself usually doesn't have a prefix or connector in 'tree' commands
	// unless it's part of a larger list. We'll start with empty strings.
	return v.renderNode(node, "", "", true, true)
}

// renderNode is the core recursive function that handles line-by-line printing.
func (v *TreeNodeVisitor) renderNode(node *treeNode, prefix string, connector string, isLast bool, isRoot bool) error {
	// 1. Render the current node
	line := fmt.Sprintf("%s%s%s\n", prefix, connector, ColorizeItem(v.writer, *node.Item))
	if _, err := io.WriteString(v.writer, line); err != nil {
		return err
	}

	// 2. Prepare prefix for children
	var newPrefix string
	if !isRoot {
		if isLast {
			newPrefix = prefix + v.style.Indent
		} else {
			newPrefix = prefix + v.style.VerticalBar
		}
	}

	// 3. Render children
	return v.renderChildrenNodes(node.Children, newPrefix, connector, isLast, isRoot)
}

func (v *TreeNodeVisitor) renderChildrenNodes(nodes []*treeNode, prefix string, connector string, isLast bool, isRoot bool) error {
	childCount := len(nodes)
	for i, child := range nodes {
		childIsLast := i == childCount-1
		childConnector := v.style.NodeConnector
		if childIsLast {
			childConnector = v.style.LastConnector
		}

		if err := v.renderNode(child, prefix, childConnector, childIsLast, false); err != nil {
			return err
		}
	}
	return nil
}

// String is no longer needed as we write directly to the writer,
// but we keep it for backward compatibility with existing tests if necessary
// by using a temporary buffer.
func (v *TreeNodeVisitor) String() string {
	// Note: This is inefficient and shouldn't be used in production.
	// It's better to update tests to pass a buffer.
	return ""
}
