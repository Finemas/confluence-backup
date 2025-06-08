// main.go
package main

import (
	"fmt"
	"log"
)

func main() {
	fmt.Println("🚀 Confluence backup tool started")

	config := LoadConfig()

	pages, err := FetchAllPages(config)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("All pages(%d)\n", len(pages))

	currentPages := FilterPages(pages, func(p Page) bool {
		return p.Status == "current"
	})
	fmt.Printf("Current pages(%d)\n", len(currentPages))

	rootPages := FilterPages(currentPages, func(p Page) bool {
		return len(p.Ancestors) == 1
	})
	fmt.Printf("Root pages(%d)\n", len(rootPages))

	// for _, page := range rootPages {
	// 	fmt.Println("-", page.Title, page.ID)

	// 	if strings.Contains(page.Title, "Mobile DevOps") {
	// 		PrintChildrenRecursive(page.ID, "\t", config)
	// 	}
	// }

	pageID := "3478781981"
	jsonErr := PrintPageJSONByID(pageID, config)
	if jsonErr != nil {
		fmt.Println("OOPPPS")
	}

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

func PrintChildrenRecursive(parentID string, indent string, cfg Config) error {
	children, err := FetchChildPages(parentID, cfg)
	if err != nil {
		return err
	}

	for _, child := range children {
		fmt.Printf("%s- %s [%s]\n", indent, child.Title, child.ID)
		if err := PrintChildrenRecursive(child.ID, indent+"  ", cfg); err != nil {
			return err
		}
	}

	return nil
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
