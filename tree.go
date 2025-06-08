package main

import (
	"fmt"
	"sort"
	"strings"
)

func BuildPageTree(pages []Page) *PageNode {
	idToNode := make(map[string]*PageNode)
	childIDSet := make(map[string]bool)

	// Build all nodes first
	for _, p := range pages {
		// Optional: skip empty or suspicious titles
		if strings.TrimSpace(p.Title) == "" || strings.HasPrefix(p.Title, "att") {
			continue
		}
		idToNode[p.ID] = &PageNode{Page: p}
	}

	// Link children and record them
	for _, p := range pages {
		node := idToNode[p.ID]
		if node == nil {
			continue
		}
		if len(p.Ancestors) == 0 {
			continue // maybe root, we'll check later
		}

		parentID := p.Ancestors[len(p.Ancestors)-1].ID
		parent := idToNode[parentID]
		if parent != nil {
			parent.Children = append(parent.Children, node)
			childIDSet[p.ID] = true
		}
	}

	// Detect true root(s) — not referenced as children
	root := &PageNode{Page: Page{Title: pages[0].Space.Key + " Root", ID: "root"}}
	for id, node := range idToNode {
		if !childIDSet[id] {
			root.Children = append(root.Children, node)
		}
	}

	return root
}

// // PrintTree recursively prints the tree with indentation.
func PrintTree(node *PageNode, indent string) {
	fmt.Printf("%s- %s(%s) [%s - %s]\n", indent, node.Page.Title, node.Page.Type, node.Page.ID, node.Page.Status)
	for _, child := range node.Children {
		PrintTree(child, indent+"  ")
	}
}

func sortNodes(nodes []*PageNode) {
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Page.Title < nodes[j].Page.Title
	})
}
