package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func fetchSpaceHome(cfg Config) (*Page, error) {
	apiURL := fmt.Sprintf("%s/rest/api/space/%s?expand=homepage", cfg.BaseURL, cfg.SpaceKey)
	req, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(cfg.Email, cfg.Token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch space info: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("space fetch HTTP %d", resp.StatusCode)
	}

	var spaceData struct {
		Homepage Page `json:"homepage"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&spaceData); err != nil {
		return nil, fmt.Errorf("parse space homepage: %w", err)
	}
	return &spaceData.Homepage, nil
}

func fetchAllChildren(cfg Config, parent Page) (*PageNode, error) {
	node := &PageNode{Page: parent, Children: []*PageNode{}}
	start := 0

	for {
		childURL := fmt.Sprintf("%s/rest/api/content/%s/child/page?limit=100&start=%d", cfg.BaseURL, parent.ID, start)
		req, err := http.NewRequest(http.MethodGet, childURL, nil)
		if err != nil {
			return nil, err
		}
		req.SetBasicAuth(cfg.Email, cfg.Token)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch children of page %s: %w", parent.ID, err)
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("children fetch HTTP %d for page %s", resp.StatusCode, parent.ID)
		}

		var result PageResponse
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode children of %s: %w", parent.ID, err)
		}

		for _, child := range result.Results {
			if child.Status != "current" {
				continue
			}
			childNode, err := fetchAllChildren(cfg, child)
			if err != nil {
				log.Printf("Warning: failed to fetch children of %s: %v", child.ID, err)
				continue
			}
			node.Children = append(node.Children, childNode)
		}

		if result.Links.Next == "" {
			break
		}
		start += len(result.Results)
	}
	return node, nil
}

func FetchSpaceHierarchy(cfg Config) (*PageNode, error) {
	rootPage, err := fetchSpaceHome(cfg)
	if err != nil {
		return nil, err
	}
	return fetchAllChildren(cfg, *rootPage)
}

func fetchAllPagesByCQL(cfg Config) ([]Page, error) {
	var allPages []Page

	// Construct the initial request URL using CQL
	cql := fmt.Sprintf("space=\"%s\" AND type IN (page, folder)", cfg.SpaceKey)
	cqlContext := `{"spaceKey":"` + cfg.SpaceKey + `","contentStatuses":["current","archived"]}`

	params := url.Values{}
	params.Set("cql", cql)
	params.Set("cqlcontext", cqlContext)
	params.Set("limit", "100")
	params.Set("expand", "ancestors")

	url := fmt.Sprintf("%s/rest/api/content/search?%s", cfg.BaseURL, params.Encode())

	for url != "" {
		req, err := BuildRequest(url, cfg)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		data, err := DoRequest(req)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		var result PageResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		allPages = append(allPages, result.Results...)

		// Print debug info (optional)
		fmt.Println("Next link:", result.Links.Next)
		fmt.Println("Fetched pages:", len(result.Results))

		// Follow the cursor-based pagination
		if result.Links.Next != "" {
			url = cfg.BaseURL + result.Links.Next
		} else {
			url = ""
		}
	}

	return allPages, nil
}

// func FetchAllPages(cfg Config) ([]Page, error) {
// 	var allPages []Page
// 	start := 0

// 	for {
// 		url := fmt.Sprintf(
// 			"%s/rest/api/content?type=page&spaceKey=%s&limit=100&start=%d&expand=ancestors,body.storage,metadata.properties.archived",
// 			cfg.BaseURL,
// 			cfg.SpaceKey,
// 			start,
// 		)

// 		req, err := BuildRequest(url, cfg)
// 		if err != nil {
// 			return nil, err
// 		}

// 		data, err := DoRequest(req)
// 		if err != nil {
// 			return nil, fmt.Errorf("read response: %w", err)
// 		}

// 		var result PageResponse
// 		if err := json.Unmarshal(data, &result); err != nil {
// 			return nil, err
// 		}

// 		allPages = append(allPages, result.Results...)

// 		if result.Links.Next == "" {
// 			break
// 		}
// 		start += len(result.Results)
// 	}

// 	return allPages, nil
// }

// func PrintPageJSONByID(pageID string, cfg Config) error {
// 	expand := url.QueryEscape("body.storage,ancestors,metadata.properties.archived,space,version")
// 	url := fmt.Sprintf("%s/rest/api/content/%s?spaceKey=%s&expand=%s", cfg.BaseURL, pageID, cfg.SpaceKey, expand)

// 	req, err := BuildRequest(url, cfg)
// 	if err != nil {
// 		return fmt.Errorf("build request: %w", err)
// 	}

// 	data, err := DoRequest(req)
// 	if err != nil {
// 		return fmt.Errorf("read response: %w", err)
// 	}

// 	var pretty map[string]interface{}
// 	if err := json.Unmarshal(data, &pretty); err != nil {
// 		return fmt.Errorf("unmarshal json: %w", err)
// 	}

// 	output, _ := json.MarshalIndent(pretty, "", "  ")
// 	fmt.Println(string(output))
// 	return nil
// }

// func FetchChildPages(parentID string, cfg Config) ([]Page, error) {
// 	var allChildren []Page
// 	start := 0

// 	for {
// 		url := fmt.Sprintf(
// 			"%s/rest/api/content/%s/child/page?limit=100&start=%d&expand=body.storage,metadata.properties.archived",
// 			cfg.BaseURL,
// 			parentID,
// 			start,
// 		)

// 		req, err := BuildRequest(url, cfg)
// 		if err != nil {
// 			return nil, fmt.Errorf("build request: %w", err)
// 		}

// 		data, err := DoRequest(req)
// 		if err != nil {
// 			return nil, fmt.Errorf("do request: %w", err)
// 		}

// 		var result PageResponse
// 		if err := json.Unmarshal(data, &result); err != nil {
// 			return nil, fmt.Errorf("unmarshal response: %w", err)
// 		}

// 		// Defensive: ensure children have valid IDs and titles
// 		for _, page := range result.Results {
// 			if page.ID != "" && page.Title != "" {
// 				allChildren = append(allChildren, page)
// 			}
// 		}

// 		if result.Links.Next == "" {
// 			break
// 		}
// 		start += len(result.Results)
// 	}

// 	return allChildren, nil
// }

// func GetPageIDByTitleInSpace(title string, cfg Config) (string, error) {
// 	cql := fmt.Sprintf(`type=page AND space="%s" AND title ~ "%s"`, cfg.SpaceKey, title)
// 	query := url.QueryEscape(cql)
// 	url := fmt.Sprintf("%s/rest/api/content/search?cql=%s", cfg.BaseURL, query)

// 	req, _ := BuildRequest(url, cfg)
// 	resp, _ := http.DefaultClient.Do(req)
// 	defer resp.Body.Close()

// 	data, _ := io.ReadAll(resp.Body)
// 	var result PageResponse
// 	json.Unmarshal(data, &result)

// 	// fmt.Println(result)
// 	// fmt.Println(result.Links)
// 	fmt.Println(result.Results)
// 	fmt.Println(len(result.Results))
// 	// fmt.Println(result.Links.Next)
// 	if len(result.Results) == 0 {
// 		return "", fmt.Errorf("no page found for title ~ \"%s\" in space %s", title, cfg.SpaceKey)
// 	}

// 	// ✅ Return first match
// 	return result.Results[0].ID, nil
// }
