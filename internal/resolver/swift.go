package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prohv/watchdocs-cli/internal/models"
)

type githubSearchItem struct {
	HTMLURL string `json:"html_url"`
}

type githubSearchResponse struct {
	Items []githubSearchItem `json:"items"`
}

func OnlineSwiftResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.TrimSpace(dep.Name)
	if packageName == "" {
		return nil, nil
	}

	// 1. If package name looks like a URL or repository path already (e.g., github.com/foo/bar)
	if strings.Contains(packageName, "/") {
		docURL := packageName
		if !strings.HasPrefix(docURL, "http://") && !strings.HasPrefix(docURL, "https://") {
			docURL = "https://" + docURL
		}
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    docURL,
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	// 2. Query GitHub Search API for the Swift package repo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://api.github.com/search/repositories?q=%s+language:swift&per_page=1", packageName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "watchdocs-cli/1.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil // timeout or network error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Fallback to github search
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://github.com/search?q=%s+language:swift", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	var data githubSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://github.com/search?q=%s+language:swift", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	docURL := ""
	if len(data.Items) > 0 {
		docURL = data.Items[0].HTMLURL
	}

	if docURL == "" {
		docURL = fmt.Sprintf("https://github.com/search?q=%s+language:swift", packageName)
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
