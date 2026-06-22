package resolver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/prohv/watchdocs-cli/internal/models"
)

type mavenDoc struct {
	G             string `json:"g"`
	A             string `json:"a"`
	LatestVersion string `json:"latestVersion"`
	Homepage      string `json:"homepage"`
}

type mavenSearchResponse struct {
	Response *struct {
		Docs []mavenDoc `json:"docs"`
	} `json:"response"`
}

func OnlineMavenResolver(dep models.Dependency) (*models.DocResult, error) {
	parts := strings.SplitN(dep.Name, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, nil
	}
	groupID, artifactID := parts[0], parts[1]

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := fmt.Sprintf(`g:"%s" AND a:"%s"`, groupID, artifactID)
	endpoint := fmt.Sprintf(
		"https://search.maven.org/solrsearch/select?q=%s&rows=1&wt=json",
		url.QueryEscape(query),
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
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

	var data mavenSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, nil
	}

	if data.Response == nil || len(data.Response.Docs) == 0 {
		return nil, nil
	}

	doc := data.Response.Docs[0]

	docURL := doc.Homepage
	if docURL == "" {
		docURL = fmt.Sprintf("https://search.maven.org/artifact/%s/%s", doc.G, doc.A)
	}

	version := dep.Version
	if version == "" {
		version = doc.LatestVersion
	}

	return &models.DocResult{
		Name:      dep.Name,
		Version:   version,
		DocURL:    docURL,
		Ecosystem: dep.Ecosystem,
		Type:      dep.Type,
	}, nil
}
