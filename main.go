// main.go
package main

import (
	"fmt"
	"log"
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
		log.Fatal(err)
	}
	fmt.Println("All pages(%d)", len(pages))

	currentPages := FilterPages(pages, func(p Page) bool {
		return p.Status == "current"
	})
	fmt.Printf("Current pages(%d)\n", len(currentPages))

	rootPages := FilterPages(currentPages, func(p Page) bool {
		return len(p.Ancestors) == 1
	})
	fmt.Printf("Root pages(%d)\n", len(rootPages))

	for _, page := range rootPages {
		fmt.Println("-", page.Title, page.ID)
	}

	// jsonErr := PrintPageJSONByID(baseURL, email, token, spaceKey, "3325100062")
	// if jsonErr != nil {
	// 	log.Fatal(jsonErr)
	// }

	// pageTitle := "Mobile DevOps"
	// pageID, pageErr := GetPageIDByTitleInSpace(baseURL, email, token, spaceKey, pageTitle)
	// if pageErr != nil {
	// 	log.Fatal(pageErr)
	// }
	// fmt.Println(pageID)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	// children, err2 := FetchChildPages(baseURL, email, token, "4527652884")
	// if err2 != nil {
	// 	log.Fatal(err2)
	// }
	// for _, c := range children {
	// 	fmt.Printf("📄 %s (%s)\n", c.Title, c.ID)
	// }
	// pages, err := FetchAllPages(baseURL, email, token, spaceKey)
	// if err != nil {
	// 	panic(fmt.Errorf("❌ Fetch failed: %w", err))
	// }

	// fmt.Printf("📚 Total pages fetched: %d\n\n", len(pages))
	// PrintPagesTree(pages)

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
