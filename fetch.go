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
	var allPages []Page
	start := 0

	for {
		url := fmt.Sprintf(
			"%s/rest/api/content?type=page&spaceKey=%s&limit=100&start=%d&expand=ancestors,body.storage,metadata.properties.archived",
			baseURL,
			spaceKey,
			start,
		)

		req, err := BuildRequest(url, email, token)
		if err != nil {
			return nil, err
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		data, _ := io.ReadAll(resp.Body)

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

func PrintPageJSONByID(baseURL, email, token, spacekey, pageID string) error {
	expand := url.QueryEscape("body.storage,ancestors,metadata.properties.archived,space,version")
	url := fmt.Sprintf("%s/rest/api/content/%s?spaceKey=%s&expand=%s", baseURL, pageID, spacekey, expand)

	req, err := BuildRequest(url, email, token)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	var pretty map[string]interface{}
	if err := json.Unmarshal(data, &pretty); err != nil {
		return fmt.Errorf("unmarshal json: %w", err)
	}

	output, _ := json.MarshalIndent(pretty, "", "  ")
	fmt.Println(string(output))
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
