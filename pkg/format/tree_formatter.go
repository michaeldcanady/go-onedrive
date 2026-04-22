package formatting

import (
	"fmt"
	"io"
	"sort"
	"strings"

	shared "github.com/michaeldcanady/go-onedrive/internal/features/fs"
)

// TreeFormatter implements OutputFormatter to render items in a hierarchical tree layout.
type TreeFormatter struct{}

// NewTreeFormatter initializes a new instance of the TreeFormatter.
func NewTreeFormatter() *TreeFormatter {
	return &TreeFormatter{}
}

// Format constructs a tree structure from the provided items and writes the rendered visualization to the writer.
func (f *TreeFormatter) Format(w io.Writer, items []any) error {
	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "(empty)")
		return err
	}

	typedItems := make([]shared.Item, len(items))
	for i, it := range items {
		typedItems[i] = it.(shared.Item)
	}

	root := f.buildTree(typedItems)
	visitor := NewTreeNodeVisitor(w)
	return visitor.VisitNode(root)
}

func (f *TreeFormatter) buildTree(items []shared.Item) *treeNode {
	nodes := make(map[string]*treeNode)

	for i := range items {
		item := items[i]
		fullPath := strings.TrimSuffix(item.Path, "/")
		if fullPath == "" {
			fullPath = item.Name
		} else {
			fullPath = fullPath + "/" + item.Name
		}

		nodes[fullPath] = &treeNode{
			Name:     item.Name,
			Item:     &item,
			Children: []*treeNode{},
		}
	}

	allNodes := make(map[string]*treeNode)
	for fullPath, node := range nodes {
		allNodes[fullPath] = node
		parentPath := strings.TrimSuffix(node.Item.Path, "/")
		for parentPath != "" {
			if _, exists := allNodes[parentPath]; exists {
				break
			}

			segments := strings.Split(parentPath, "/")
			name := segments[len(segments)-1]
			parentNode := &treeNode{
				Name: name,
				Item: &shared.Item{
					Name: name,
					Path: strings.Join(segments[:len(segments)-1], "/"),
					Type: shared.TypeFolder,
				},
				Children: []*treeNode{},
			}
			allNodes[parentPath] = parentNode
			parentPath = strings.Join(segments[:len(segments)-1], "/")
		}
	}

	hasParent := make(map[string]bool)
	for fullPath, node := range allNodes {
		item := node.Item
		parentPath := strings.TrimSuffix(item.Path, "/")
		if parentPath == "" {
			continue
		}
		if parent, ok := allNodes[parentPath]; ok {
			parent.Children = append(parent.Children, node)
			hasParent[fullPath] = true
		}
	}

	var roots []*treeNode
	for fullPath, node := range allNodes {
		if !hasParent[fullPath] {
			roots = append(roots, node)
		}
	}

	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Name < roots[j].Name
	})

	for _, r := range roots {
		f.sortTree(r)
	}

	if len(roots) == 1 {
		return roots[0]
	}

	return &treeNode{
		Name: "/",
		Item: &shared.Item{
			Name: "/",
			Path: "",
			Type: shared.TypeFolder,
		},
		Children: roots,
	}
}

func (f *TreeFormatter) sortTree(n *treeNode) {
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
	for _, c := range n.Children {
		f.sortTree(c)
	}
}
