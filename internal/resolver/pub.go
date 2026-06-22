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

type pubLatest struct {
	Version string `json:"version"`
}

type pubResponse struct {
	Latest *pubLatest `json:"latest"`
}

func OnlinePubResolver(dep models.Dependency) (*models.DocResult, error) {
	packageName := strings.TrimSpace(dep.Name)
	if packageName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://pub.dev/api/packages/%s", packageName)
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

	var data pubResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	version := dep.Version
	if version == "" && data.Latest != nil {
		version = data.Latest.Version
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   version,
		DocURL:    fmt.Sprintf("https://pub.dev/packages/%s", packageName),
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
