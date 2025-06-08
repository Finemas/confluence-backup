package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func BuildRequest(url string, cfg Config) (*http.Request, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	if strings.Contains(url, "/api/v2/") {
		req.Header.Set("Authorization", "Bearer "+cfg.Token)
	} else {
		req.SetBasicAuth(cfg.Email, cfg.Token)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func DoRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	data, readErr := io.ReadAll(resp.Body)
	if readErr != nil {
		return nil, fmt.Errorf("read body: %w", readErr)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func PrintAsJSON(data []byte) {
	var prettyJSON map[string]interface{}
	if err := json.Unmarshal(data, &prettyJSON); err != nil {
		fmt.Printf("parse json error: %w\n", err)
		return
	}

	out, _ := json.MarshalIndent(prettyJSON, "", "  ")
	fmt.Println(string(out))
}
