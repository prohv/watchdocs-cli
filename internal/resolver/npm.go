package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/prohv/watchdocs-cli/internal/models"
)

type npmRegistryResponse struct {
	Homepage   string `json:"homepage"`
	Repository struct {
		URL string `json:"url"`
	} `json:"repository"`
}

func OnlineNpmResolver(dep models.Dependency) (*models.DocResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", dep.Name)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, nil // timeout or network error — return nil like the TS version
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}

	var data npmRegistryResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	docURL := data.Homepage
	if docURL == "" {
		docURL = data.Repository.URL
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
