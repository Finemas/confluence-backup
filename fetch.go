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

		var parsed PageResponse
		if err := json.Unmarshal(data, &parsed); err != nil {
			return nil, err
		}

		allPages = append(allPages, parsed.Results...)

		if parsed.Links.Next != "" {
			nextURL = baseURL + parsed.Links.Next
		} else {
			nextURL = ""
		}
	}

	return allPages, nil
}
