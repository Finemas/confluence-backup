// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("🚀 Confluence backup tool started")

	config := LoadConfig()

	// Step 1: Fetch all pages
	pages, err := fetchAllPagesByCQL(config)
	if err != nil {
		log.Fatalf("fetch error: %v", err)
	}
	fmt.Printf("Fetched %d pages\n", len(pages))

	// Step 2: Build tree structure from flat list
	tree := BuildPageTree(pages)

	// Step 3: Print tree to terminal
	fmt.Println("📁 Page Tree:")
	PrintTree(tree, "")

	// trees := buildTreeFromAncestors(pages)
	// for _, root := range trees {
	// 	printTree(root, "")
	// }

	// pages, err := FetchAllPages(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("All pages(%d)\n", len(pages))

	// currentPages := FilterPages(pages, func(p Page) bool {
	// 	return p.Status == "current"
	// })
	// fmt.Printf("Current pages(%d)\n", len(currentPages))

	// rootPages := FilterPages(currentPages, func(p Page) bool {
	// 	return len(p.Ancestors) == 1
	// })
	// fmt.Printf("Root pages(%d)\n", len(rootPages))

	// for _, page := range rootPages {
	// 	fmt.Println("-", page.Title, page.ID)

	// 	if strings.Contains(page.Title, "Mobile DevOps") {
	// 		PrintChildrenRecursive(page.ID, "\t", config)
	// 	}
	// }

	// pageID := "3478781981"
	// jsonErr := PrintPageJSONByID(pageID, config)
	// if jsonErr != nil {
	// 	fmt.Println("OOPPPS")
	// }

	// pageTitle := "Mobile DevOps"
	// pageID, pageErr := GetPageIDByTitleInSpace(pageTitle, config)
	// if pageErr != nil {
	// 	log.Fatal(pageErr)
	// }
	// fmt.Println(pageID)

	// // 1. Fetch all structured pages
	// pages, _ := FetchAllPages(baseURL, email, token, spaceKey)
	// for _, p := range pages {
	// 	path := []string{}
	// 	for _, a := range p.Ancestors {
	// 		path = append(path, a.Title)
	// 	}
	// 	path = append(path, p.Title)
	// 	fmt.Println("📄", strings.Join(path, " / "))
	// }
	// fmt.Println("📚 Full Page Tree")
	// PrintPagesTree(structuredPages)

	// 2. Fetch only my pages (flat)
	// myPages, _ := FetchMyPages(baseURL, email, token, spaceKey)
	// fmt.Printf("\n📄 My Pages (%d total):\n", len(myPages))
	// for _, p := range myPages {
	// 	fmt.Printf("• %s (%s)\n", p.Title, p.ID)
	// }
}

func FilterPages(pages []Page, match func(Page) bool) []Page {
	var result []Page
	for _, p := range pages {
		if match(p) {
			result = append(result, p)
		}
	}
	return result
}
