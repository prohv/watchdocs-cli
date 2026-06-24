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

type nugetPackageData struct {
	ID         string `json:"id"`
	ProjectURL string `json:"projectUrl"`
}

type nugetSearchResponse struct {
	Data []nugetPackageData `json:"data"`
}

func OnlineNuGetResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.TrimSpace(dep.Name)
	if packageName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Query NuGet Search API with optimized parameters: take=1 limits response array size
	url := fmt.Sprintf("https://azuresearch-usnc.nuget.org/query?q=packageid:%s&prerelease=false&take=1", packageName)
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
		// Fallback to official packages page
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://www.nuget.org/packages/%s", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	var data nugetSearchResponse
	// Optimization: json.NewDecoder will parse the JSON stream, and our selective struct
	// will skip deserialization of large fields like version lists, descriptions, and tag lists.
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return &models.DocResult{
			Name:      dep.Name,
			Version:   dep.Version,
			DocURL:    fmt.Sprintf("https://www.nuget.org/packages/%s", packageName),
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
		}, nil
	}

	docURL := ""
	if len(data.Data) > 0 {
		docURL = data.Data[0].ProjectURL
	}

	if docURL == "" {
		docURL = fmt.Sprintf("https://www.nuget.org/packages/%s", packageName)
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
