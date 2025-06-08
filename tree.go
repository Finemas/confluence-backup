package main

import (
	"fmt"
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

// // FlattenTree flattens the tree into a list of all pages (for saving/export).
// func FlattenTree(node *PageNode) []Page {
// 	var pages []Page
// 	var walk func(n *PageNode)

// 	walk = func(n *PageNode) {
// 		pages = append(pages, n.Page)
// 		for _, c := range n.Children {
// 			walk(c)
// 		}
// 	}
// 	walk(node)
// 	return pages
// }

// // FindPageNode searches for a page node by title or ID.
// func FindPageNode(node *PageNode, match func(Page) bool) *PageNode {
// 	if match(node.Page) {
// 		return node
// 	}
// 	for _, child := range node.Children {
// 		if found := FindPageNode(child, match); found != nil {
// 			return found
// 		}
// 	}
// 	return nil
// }

// // TreeToMarkdown returns a Markdown-formatted representation of the tree.
// func TreeToMarkdown(node *PageNode, level int) string {
// 	var sb strings.Builder
// 	prefix := strings.Repeat("#", level+1)
// 	sb.WriteString(fmt.Sprintf("%s %s\n\n", prefix, node.Page.Title))
// 	sb.WriteString(fmt.Sprintf("_Page ID: %s_\n\n", node.Page.ID))
// 	// you could also embed content, e.g., node.Page.Body.Storage.Value here

// 	for _, child := range node.Children {
// 		sb.WriteString(TreeToMarkdown(child, level+1))
// 	}
// 	return sb.String()
// }

// func buildTreeFromAncestors(pages []Page) []*Page {
// 	pageMap := map[string]*Page{}
// 	var roots []*Page
// 	for i := range pages {
// 		page := pages[i] // create local pointer
// 		page.Children = []*Page{}
// 		pageMap[page.ID] = &page
// 	}
// 	for i := range pages {
// 		page := pageMap[pages[i].ID]
// 		if len(page.Ancestors) == 0 {
// 			roots = append(roots, page)
// 		} else {
// 			parentID := page.Ancestors[len(page.Ancestors)-1].ID
// 			parent, ok := pageMap[parentID]
// 			if ok {
// 				parent.Children = append(parent.Children, page)
// 			} else {
// 				roots = append(roots, page)
// 			}
// 		}
// 	}
// 	return roots
// }

// func printTree(page *Page, indent string) {
// 	icon := "📄"
// 	if page.Status == "archived" {
// 		icon = "📦"
// 	}
// 	fmt.Printf("%s%s %s [%s]\n", indent, icon, page.Title, page.ID)
// 	sort.Slice(page.Children, func(i, j int) bool {
// 		return strings.ToLower(page.Children[i].Title) < strings.ToLower(page.Children[j].Title)
// 	})
// 	for _, child := range page.Children {
// 		printTree(child, indent+"  ")
// 	}
// }

// func PrintPagesTree(pages []Page) {
// 	spaceRoot := detectSpaceRootTitle(pages)
// 	fmt.Println("🧠 Detected space root:", spaceRoot)

// 	for _, p := range pages {
// 		fmt.Printf("🧩 %s (%s)\n", p.Title, p.ID)
// 		for i, a := range p.Ancestors {
// 			fmt.Printf("  %d. %s (%s)\n", i+1, a.Title, a.ID)
// 		}
// 		fmt.Println()
// 	}

// 	// idToNode := make(map[string]*PageNode)
// 	// for _, p := range pages {
// 	// 	idToNode[p.ID] = &PageNode{Page: p}
// 	// }

// 	// var roots []*PageNode

// 	// for _, node := range idToNode {
// 	// 	parentID := parentAncestorID(node.Page, spaceRoot)
// 	// 	if parentID == "" {
// 	// 		roots = append(roots, node)
// 	// 	} else if parent, exists := idToNode[parentID]; exists {
// 	// 		parent.Children = append(parent.Children, node)
// 	// 	} else {
// 	// 		fmt.Printf("⚠️  Orphaned page: %s (missing parent %s)\n", node.Page.Title, parentID)
// 	// 		roots = append(roots, node)
// 	// 	}
// 	// }

// 	// sortNodes(roots)
// 	// for _, root := range roots {
// 	// 	printTreeAsPaths(root, 0)
// 	// }
// }

// func printTreeAsPaths(node *PageNode, indent int) {
// 	prefix := strings.Repeat("  ", indent)

// 	if len(node.Children) > 0 {
// 		fmt.Printf("%s📁 %s/\n", prefix, node.Page.Title)
// 	} else {
// 		fmt.Printf("%s📄 %s.md\n", prefix, node.Page.Title)
// 	}

// 	sortNodes(node.Children)
// 	for _, child := range node.Children {
// 		printTreeAsPaths(child, indent+1)
// 	}
// }

// // --- Ancestor helpers ---

// func detectSpaceRootTitle(pages []Page) string {
// 	count := make(map[string]int)

// 	for _, p := range pages {
// 		if len(p.Ancestors) > 0 {
// 			root := p.Ancestors[0].Title
// 			count[root]++
// 		}
// 	}

// 	var mostCommon string
// 	max := 0
// 	for title, c := range count {
// 		if c > max {
// 			mostCommon = title
// 			max = c
// 		}
// 	}

// 	return mostCommon
// }

// func cleanAncestors(p Page, spaceRoot string) []Ancestor {
// 	anc := p.Ancestors
// 	if len(anc) > 0 && anc[0].Title == spaceRoot {
// 		return anc[1:]
// 	}
// 	return anc
// }

// func parentAncestorID(p Page, spaceRoot string) string {
// 	ancestors := cleanAncestors(p, spaceRoot)
// 	if len(ancestors) == 0 {
// 		return ""
// 	}
// 	return ancestors[len(ancestors)-1].ID
// }

// func sortNodes(nodes []*PageNode) {
// 	sort.SliceStable(nodes, func(i, j int) bool {
// 		return nodes[i].Page.Title < nodes[j].Page.Title
// 	})
// }
