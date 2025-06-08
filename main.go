// main.go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("🚀 Confluence backup tool started")

	LoadEnv() // from config.go

	baseURL := GetEnv("DRMAX_CONFLUENCE_BASE")
	email := GetEnv("DRMAX_CONFLUENCE_EMAIL")
	token := GetEnv("DRMAX_CONFLUENCE_TOKEN")
	spaceKey := GetEnv("DRMAX_SPACE_KEY")

	pages, err := FetchAllPages(baseURL, email, token, spaceKey)
	if err != nil {
		panic(fmt.Errorf("❌ Fetch failed: %w", err))
	}

	fmt.Printf("📚 Total pages fetched: %d\n\n", len(pages))
	PrintPagesTree(pages)

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
