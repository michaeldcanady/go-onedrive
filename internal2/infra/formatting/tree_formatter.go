package formatting

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
	domainfs "github.com/michaeldcanady/go-onedrive/internal2/domain/fs"
)

type TreeFormatter struct{}

func NewTreeFormatter() *TreeFormatter {
	return &TreeFormatter{}
}

func (f *TreeFormatter) Format(w io.Writer, v any) error {
	// Accept both []domainfs.Item and []*domainfs.Item
	items, ok := v.([]domainfs.Item)
	if !ok {
		return fmt.Errorf("unsupported element type %T", v)
	}

	if len(items) == 0 {
		_, err := fmt.Fprintln(w, "(empty)")
		return err
	}

	root := f.buildTree(items)

	visitor := &TreeNodeVisitor{}

	if err := visitor.VisitNode(root); err != nil {
		return err
	}
	_, err := w.Write([]byte(visitor.String() + "\n"))
	return err
}

func (f *TreeFormatter) buildTree(items []domainfs.Item) *treeNode {
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
				Item: &domainfs.Item{
					Name: segments[len(segments)-1],
					Path: strings.Join(segments[:len(segments)-1], "/"),
					Type: fs.ItemTypeFolder,
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
	for _, r := range roots {
		sortTree(r)
	}

	// 5. Single root or synthetic root
	if len(roots) == 1 {
		return roots[0]
	}

	return &treeNode{
		Name: "/",
		Item: &domainfs.Item{
			Name: "/",
			Path: "",
			Type: fs.ItemTypeFolder,
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
