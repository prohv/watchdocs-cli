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

type packagistPackage struct {
	Repository string `json:"repository"`
	Homepage   string `json:"homepage"`
}

type packagistResponse struct {
	Package packagistPackage `json:"package"`
}

func OnlineComposerResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.TrimSpace(dep.Name)
	if packageName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query Packagist API
	url := fmt.Sprintf("https://packagist.org/packages/%s.json", packageName)
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
		// Fallback to packagist listing page
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://packagist.org/packages/%s", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	var data packagistResponse
	// Optimization: Selective struct decoding ignores the huge version history JSON list,
	// only reading the repository and homepage fields at top-level.
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://packagist.org/packages/%s", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	docURL := data.Package.Homepage
	if docURL == "" {
		docURL = data.Package.Repository
	}
	if docURL == "" {
		docURL = fmt.Sprintf("https://packagist.org/packages/%s", packageName)
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
