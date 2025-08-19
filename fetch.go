package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"
)

func FetchPageByID(pageID string, cfg Config) (*Page, error) {
	expand := url.QueryEscape("body.storage,ancestors,metadata.properties.archived,space,version")
	url := fmt.Sprintf("%s/%s/%s?spaceKey=%s&expand=%s", cfg.BaseURL, AllPages, pageID, cfg.SpaceKey, expand)

	req, err := BuildRequest(url, cfg)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	data, err := DoRequest(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	var page Page
	if err := json.Unmarshal(data, &page); err != nil {
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return &page, nil
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

		// Follow the cursor-based pagination
		if result.Links.Next != "" {
			url = cfg.BaseURL + result.Links.Next
		} else {
			url = ""
		}
	}

	return allPages, nil
}

func wasContributedBy(cfg Config, pageID, userAccountID string) (bool, error) {
	url := fmt.Sprintf("%s/%s/%s/history", cfg.BaseURL, AllPages, pageID)
	req, err := BuildRequest(url, cfg)
	if err != nil {
		return false, err
	}

	data, err := DoRequest(req)
	if err != nil {
		return false, err
	}

	var history struct {
		LastUpdated struct {
			By struct {
				AccountID string `json:"accountId"`
			} `json:"by"`
		} `json:"lastUpdated"`
	}

	if err := json.Unmarshal(data, &history); err != nil {
		return false, err
	}

	return history.LastUpdated.By.AccountID == userAccountID, nil
}

func FetchContributedPages(cfg Config, allPages []Page, userAccountID string) ([]Page, error) {
	var result []Page
	for _, page := range allPages {
		ok, err := wasContributedBy(cfg, page.ID, userAccountID)
		if err != nil {
			log.Printf("Skipping %s: %v", page.Title, err)
			continue
		}
		if ok {
			result = append(result, page)
		}
		time.Sleep(150 * time.Millisecond)
	}
	return result, nil
}

func PrintPageJSONByID(pageID string, cfg Config) error {
	expand := url.QueryEscape("body.storage,ancestors,metadata.properties.archived,space,version")
	url := fmt.Sprintf("%s/%s/%s?spaceKey=%s&expand=%s", cfg.BaseURL, AllPages, pageID, cfg.SpaceKey, expand)

	req, err := BuildRequest(url, cfg)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	data, err := DoRequest(req)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	var pretty map[string]interface{}
	if err := json.Unmarshal(data, &pretty); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	output, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Println(string(output))
	return nil
}

func GetPageIDByTitleInSpace(title string, cfg Config) (string, error) {
	cql := fmt.Sprintf(`type=page AND space="%s" AND title ~ "%s"`, cfg.SpaceKey, title)
	query := url.QueryEscape(cql)
	url := fmt.Sprintf("%s/%s?cql=%s", cfg.BaseURL, SearchPages, query)

	req, _ := BuildRequest(url, cfg)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var result PageResponse
	json.Unmarshal(data, &result)

	// fmt.Println(result)
	// fmt.Println(result.Links)
	fmt.Println(result.Results)
	fmt.Println(len(result.Results))
	// fmt.Println(result.Links.Next)
	if len(result.Results) == 0 {
		return "", fmt.Errorf("no page found for title ~ \"%s\" in space %s", title, cfg.SpaceKey)
	}

	// ✅ Return first match
	return result.Results[0].ID, nil
}
