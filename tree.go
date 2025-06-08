package main

import (
	"fmt"
	"sort"
	"strings"
)

func PrintPagesTree(pages []Page) {
	spaceRoot := detectSpaceRootTitle(pages)
	fmt.Println("🧠 Detected space root:", spaceRoot)

	for _, p := range pages {
		fmt.Printf("🧩 %s (%s)\n", p.Title, p.ID)
		for i, a := range p.Ancestors {
			fmt.Printf("  %d. %s (%s)\n", i+1, a.Title, a.ID)
		}
		fmt.Println()
	}

	// idToNode := make(map[string]*PageNode)
	// for _, p := range pages {
	// 	idToNode[p.ID] = &PageNode{Page: p}
	// }

	// var roots []*PageNode

	// for _, node := range idToNode {
	// 	parentID := parentAncestorID(node.Page, spaceRoot)
	// 	if parentID == "" {
	// 		roots = append(roots, node)
	// 	} else if parent, exists := idToNode[parentID]; exists {
	// 		parent.Children = append(parent.Children, node)
	// 	} else {
	// 		fmt.Printf("⚠️  Orphaned page: %s (missing parent %s)\n", node.Page.Title, parentID)
	// 		roots = append(roots, node)
	// 	}
	// }

	// sortNodes(roots)
	// for _, root := range roots {
	// 	printTreeAsPaths(root, 0)
	// }
}

func printTreeAsPaths(node *PageNode, indent int) {
	prefix := strings.Repeat("  ", indent)

	if len(node.Children) > 0 {
		fmt.Printf("%s📁 %s/\n", prefix, node.Page.Title)
	} else {
		fmt.Printf("%s📄 %s.md\n", prefix, node.Page.Title)
	}

	sortNodes(node.Children)
	for _, child := range node.Children {
		printTreeAsPaths(child, indent+1)
	}
}

// --- Ancestor helpers ---

func detectSpaceRootTitle(pages []Page) string {
	count := make(map[string]int)

	for _, p := range pages {
		if len(p.Ancestors) > 0 {
			root := p.Ancestors[0].Title
			count[root]++
		}
	}

	var mostCommon string
	max := 0
	for title, c := range count {
		if c > max {
			mostCommon = title
			max = c
		}
	}

	return mostCommon
}

func cleanAncestors(p Page, spaceRoot string) []Ancestor {
	anc := p.Ancestors
	if len(anc) > 0 && anc[0].Title == spaceRoot {
		return anc[1:]
	}
	return anc
}

func parentAncestorID(p Page, spaceRoot string) string {
	ancestors := cleanAncestors(p, spaceRoot)
	if len(ancestors) == 0 {
		return ""
	}
	return ancestors[len(ancestors)-1].ID
}

func sortNodes(nodes []*PageNode) {
	sort.SliceStable(nodes, func(i, j int) bool {
		return nodes[i].Page.Title < nodes[j].Page.Title
	})
}
