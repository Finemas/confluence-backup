// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("🚀 Confluence backup tool started")

	config := LoadConfig()

	page, err := FetchPageByID("4579164404", config)
	if err != nil {
		log.Fatal(err)
	}

	if err := SavePageAsMarkdown(page, "pages"); err != nil {
		log.Fatal(err)
	}

	// allPages, allErr := fetchAllPagesByCQL(config)
	// if allErr != nil {
	// 	log.Fatalf("fetch error: %v", allErr)
	// }
	// fmt.Printf("All pages(%d)\n", len(allPages))

	// contribPages, contrigErr := FetchContributedPages(config, allPages, "<your-account-id>")
	// if contrigErr != nil {
	// 	log.Fatalf("fetch error: %v", contrigErr)
	// }
	// fmt.Printf("Contrib pages(%d)\n", len(contribPages))

	// tree := BuildPageTree(allPages)

	// fmt.Println("📁 Tree:")
	// PrintTree(tree, "")
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
