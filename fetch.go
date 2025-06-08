package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

func BuildRequest(url, email, token string) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(email, token)
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func PrintRawJSON(url, email, token string) error {
	req, err := BuildRequest(url, email, token)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	out, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println(string(out))
	return nil
}

func FetchAllPages(baseURL, email, token, spaceKey string) ([]Page, error) {
	params := url.Values{}
	params.Set("type", "page")
	params.Set("spaceKey", spaceKey)
	params.Set("expand", "body.storage,ancestors")
	params.Set("limit", "100")

	nextURL := fmt.Sprintf("%s/rest/api/content?%s", baseURL, params.Encode())
	var allPages []Page

	for nextURL != "" {
		req, err := BuildRequest(nextURL, email, token)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var pretty map[string]interface{}
		if err := json.Unmarshal(data, &pretty); err != nil {
			return nil, fmt.Errorf("json unmarshal: %w", err)
		}

		jsonText, _ := json.MarshalIndent(pretty, "", "  ")
		fmt.Println(string(jsonText))

		// var parsed PageResponse
		// if err := json.Unmarshal(data, &parsed); err != nil {
		// 	return nil, err
		// }

		// allPages = append(allPages, parsed.Results...)

		// if parsed.Links.Next != "" {
		// 	nextURL = baseURL + parsed.Links.Next
		// } else {
		// 	nextURL = ""
		// }
	}

	return allPages, nil
}

func FetchFullPageJSON(baseURL, email, token, spaceKey, title string) error {
	// 1. Search the page by title and space
	cql := fmt.Sprintf(`type=page AND space="%s" AND title ~ "%s"`, spaceKey, title)
	query := url.QueryEscape(cql)
	searchURL := fmt.Sprintf("%s/rest/api/content/search?cql=%s", baseURL, query)

	req, err := BuildRequest(searchURL, email, token)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp.Body.Close()

	searchData, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read search response: %w", err)
	}

	var searchResult PageResponse
	if err := json.Unmarshal(searchData, &searchResult); err != nil {
		return fmt.Errorf("unmarshal search response: %w", err)
	}

	if len(searchResult.Results) == 0 {
		return fmt.Errorf("no page found for title: %s", title)
	}

	pageID := searchResult.Results[0].ID
	fmt.Println("✅ Found page ID:", pageID)

	// 2. Fetch the full page by ID
	fullURL := fmt.Sprintf("%s/rest/api/content/%s?expand=body.storage,ancestors", baseURL, pageID)
	req2, err := BuildRequest(fullURL, email, token)
	if err != nil {
		return fmt.Errorf("build full request: %w", err)
	}

	resp2, err := http.DefaultClient.Do(req2)
	if err != nil {
		return fmt.Errorf("http call failed: %w", err)
	}
	defer resp2.Body.Close()

	pageData, err := io.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("read full page body: %w", err)
	}

	var pretty map[string]interface{}
	if err := json.Unmarshal(pageData, &pretty); err != nil {
		return fmt.Errorf("unmarshal page json: %w", err)
	}

	out, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Println(string(out))
	return nil
}

func FetchChildPages(baseURL, email, token, parentID string) ([]Page, error) {
	var children []Page
	start := 0

	for {
		url := fmt.Sprintf(
			"%s/rest/api/content/%s/child/page?limit=100&start=%d&expand=body.storage,ancestors",
			baseURL,
			parentID,
			start,
		)

		req, err := BuildRequest(url, email, token)
		if err != nil {
			return nil, fmt.Errorf("build request: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("http call failed: %w", err)
		}
		defer resp.Body.Close()

		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("read response: %w", err)
		}

		var result PageResponse
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		children = append(children, result.Results...)

		if result.Links.Next == "" {
			break
		}
		start += len(result.Results)
	}

	return children, nil
}

func GetPageIDByTitleInSpace(baseURL, email, token, spaceKey, title string) (string, error) {
	cql := fmt.Sprintf(`type=page AND space="%s" AND title ~ "%s"`, spaceKey, title)
	query := url.QueryEscape(cql)
	url := fmt.Sprintf("%s/rest/api/content/search?cql=%s", baseURL, query)

	req, _ := BuildRequest(url, email, token)
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
		return "", fmt.Errorf("no page found for title ~ \"%s\" in space %s", title, spaceKey)
	}

	// ✅ Return first match
	return result.Results[0].ID, nil
}
