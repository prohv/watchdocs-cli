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

type pypiProjectURLs map[string]*string

type pypiPackageInfo struct {
	HomePage    *string         `json:"home_page"`
	ProjectURLs pypiProjectURLs `json:"project_urls"`
}

type pypiResponse struct {
	Info *pypiPackageInfo `json:"info"`
}

func OnlinePipResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.ToLower(strings.TrimSpace(dep.Name))
	if packageName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://pypi.org/pypi/%s/json", packageName)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil // timeout or network error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var data pypiResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	var docURL string
	if info := data.Info; info != nil {
		if v := info.ProjectURLs["Documentation"]; v != nil && *v != "" {
			docURL = *v
		} else if v := info.ProjectURLs["Homepage"]; v != nil && *v != "" {
			docURL = *v
		} else if info.HomePage != nil && *info.HomePage != "" {
			docURL = *info.HomePage
		}
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
