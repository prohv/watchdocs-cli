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

type crateData struct {
	Documentation *string `json:"documentation"`
	Homepage      *string `json:"homepage"`
	Repository    *string `json:"repository"`
}

type cratesIoResponse struct {
	Crate *crateData `json:"crate"`
}

func OnlineCargoResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.ToLower(strings.TrimSpace(dep.Name))
	if packageName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://crates.io/api/v1/crates/%s", packageName)
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
		return nil, nil
	}

	var data cratesIoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	var docURL string
	if c := data.Crate; c != nil {
		if c.Documentation != nil && *c.Documentation != "" {
			docURL = *c.Documentation
		} else if c.Homepage != nil && *c.Homepage != "" {
			docURL = *c.Homepage
		} else if c.Repository != nil && *c.Repository != "" {
			docURL = *c.Repository
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
