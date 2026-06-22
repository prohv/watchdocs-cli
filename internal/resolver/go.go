package resolver

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func OnlineGoResolver(dep models.Dependency) (*models.DocResult, error) {
	moduleName := strings.TrimSpace(dep.Name)
	if moduleName == "" {
		return nil, nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	url := fmt.Sprintf("https://proxy.golang.org/%s/@v/list", moduleName)
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

	var docURL string
	if dep.Version != "" {
		docURL = fmt.Sprintf("https://pkg.go.dev/%s@%s", moduleName, dep.Version)
	} else {
		docURL = fmt.Sprintf("https://pkg.go.dev/%s", moduleName)
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   dep.Version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
