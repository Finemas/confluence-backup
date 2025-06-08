package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func FetchAllPages(cfg Config) ([]Page, error) {
	var allPages []Page
	start := 0

	for {
		url := fmt.Sprintf(
			"%s/rest/api/content?type=page&spaceKey=%s&limit=100&start=%d&expand=ancestors,body.storage,metadata.properties.archived",
			cfg.BaseURL,
			cfg.SpaceKey,
			start,
		)

		req, err := BuildRequest(url, cfg)
		if err != nil {
			return nil, err
		}

		data, err := DoRequest(req)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		var result PageResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, err
		}

		allPages = append(allPages, result.Results...)

		if result.Links.Next == "" {
			break
		}
		start += len(result.Results)
	}

	return allPages, nil
}

func PrintPageJSONByID(pageID string, cfg Config) error {
	expand := url.QueryEscape("body.storage,ancestors,metadata.properties.archived,space,version")
	url := fmt.Sprintf("%s/rest/api/content/%s?spaceKey=%s&expand=%s", cfg.BaseURL, pageID, cfg.SpaceKey, expand)

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

func FetchChildPages(parentID string, cfg Config) ([]Page, error) {
	var allChildren []Page
	start := 0

	for {
		url := fmt.Sprintf(
			"%s/rest/api/content/%s/child/page?limit=100&start=%d&expand=body.storage,metadata.properties.archived",
			cfg.BaseURL,
			parentID,
			start,
		)

		req, err := BuildRequest(url, cfg)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		data, err := DoRequest(req)
		if err != nil {
			return nil, fmt.Errorf("do request: %w", err)
		}

		var result PageResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}

		// Defensive: ensure children have valid IDs and titles
		for _, page := range result.Results {
			if page.ID != "" && page.Title != "" {
				allChildren = append(allChildren, page)
			}
		}

		if result.Links.Next == "" {
			break
		}
		start += len(result.Results)
	}

	return allChildren, nil
}

func GetPageIDByTitleInSpace(title string, cfg Config) (string, error) {
	cql := fmt.Sprintf(`type=page AND space="%s" AND title ~ "%s"`, cfg.SpaceKey, title)
	query := url.QueryEscape(cql)
	url := fmt.Sprintf("%s/rest/api/content/search?cql=%s", cfg.BaseURL, query)

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
