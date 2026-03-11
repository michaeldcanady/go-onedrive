package formatting

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/fs/domain"
)

type TreeFormatter struct{}

func NewTreeFormatter() *TreeFormatter {
	return &TreeFormatter{}
}

func (f *TreeFormatter) Format(w io.Writer, items []domain.Item) error {
	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "(empty)")
		return err
	}

	root := f.buildTree(items)

	visitor := NewTreeNodeVisitor(w)

	return visitor.VisitNode(root)
}

func (f *TreeFormatter) buildTree(items []domain.Item) *treeNode {
	nodes := make(map[string]*treeNode)

	if len(items) == 0 {
		return nil
	}

	// 1. Create all nodes
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

	// 2. Build a new map including synthetic parents
	allNodes := make(map[string]*treeNode)

	for fullPath, node := range nodes {
		allNodes[fullPath] = node

		parentPath := strings.TrimSuffix(node.Item.Path, "/")
		for parentPath != "" {
			if _, exists := allNodes[parentPath]; exists {
				break
			}

			segments := strings.Split(parentPath, "/")
			parentNode := &treeNode{
				Name: segments[len(segments)-1],
				Item: &domain.Item{
					Name: segments[len(segments)-1],
					Path: strings.Join(segments[:len(segments)-1], "/"),
					Type: domain.ItemTypeFolder,
				},
				Children: []*treeNode{},
			}

			allNodes[parentPath] = parentNode

			// Walk upward
			parentPath = strings.Join(segments[:len(segments)-1], "/")
		}
	}

	// 3. Attach children
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

	// 4. Collect roots
	var roots []*treeNode
	for fullPath, node := range allNodes {
		if !hasParent[fullPath] {
			roots = append(roots, node)
		}
	}

	// 5. Sort children recursively
	sort.Slice(roots, func(i, j int) bool {
		return roots[i].Name < roots[j].Name
	})

	for _, r := range roots {
		sortTree(r)
	}

	// 5. Single root or synthetic root
	if len(roots) == 1 {
		return roots[0]
	}

	return &treeNode{
		Name: "/",
		Item: &domain.Item{
			Name: "/",
			Path: "",
			Type: domain.ItemTypeFolder,
		},
		Children: roots,
	}
}

func sortTree(n *treeNode) {
	sort.Slice(n.Children, func(i, j int) bool {
		return n.Children[i].Name < n.Children[j].Name
	})
	for _, c := range n.Children {
		sortTree(c)
	}
}
