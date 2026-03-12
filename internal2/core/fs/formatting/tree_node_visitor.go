package formatting

import (
	"fmt"
	"io"
)

// TreeStyle defines the characters used to visually represent the tree structure.
type TreeStyle struct {
	// NodeConnector is the prefix for a standard child node.
	NodeConnector string
	// LastConnector is the prefix for the final child node in a group.
	LastConnector string
	// VerticalBar is the symbol for a continuing branch line.
	VerticalBar   string
	// Indent is the spacing used for nested levels.
	Indent        string
}

// DefaultTreeStyle provides a standard set of Unicode characters for tree rendering.
var DefaultTreeStyle = TreeStyle{
	NodeConnector: "├──",
	LastConnector: "└──",
	VerticalBar:   "│   ",
	Indent:        "    ",
}

// TreeNodeVisitor orchestrates the recursive printing of a tree structure to a writer.
type TreeNodeVisitor struct {
	// writer is the destination for the rendered output.
	writer io.Writer
	// style defines the visual symbols used for connectors.
	style  TreeStyle
}

// NewTreeNodeVisitor initializes a new visitor with the default styling.
func NewTreeNodeVisitor(w io.Writer) *TreeNodeVisitor {
	return &TreeNodeVisitor{
		writer: w,
		style:  DefaultTreeStyle,
	}
}

// WithStyle configures the visitor to use a custom set of visual symbols.
func (v *TreeNodeVisitor) WithStyle(style TreeStyle) *TreeNodeVisitor {
	v.style = style
	return v
}

// VisitNode initiates the recursive rendering process from the provided root node.
func (v *TreeNodeVisitor) VisitNode(node *treeNode) error {
	if node == nil {
		return nil
	}
	return v.renderNode(node, "", "", true, true)
}

func (v *TreeNodeVisitor) renderNode(node *treeNode, prefix string, connector string, isLast bool, isRoot bool) error {
	line := fmt.Sprintf("%s%s%s\n", prefix, connector, ColorizeItem(v.writer, *node.Item))
	if _, err := io.WriteString(v.writer, line); err != nil {
		return err
	}

	var newPrefix string
	if !isRoot {
		if isLast {
			newPrefix = prefix + v.style.Indent
		} else {
			newPrefix = prefix + v.style.VerticalBar
		}
	}

	return v.renderChildrenNodes(node.Children, newPrefix)
}

func (v *TreeNodeVisitor) renderChildrenNodes(nodes []*treeNode, prefix string) error {
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
