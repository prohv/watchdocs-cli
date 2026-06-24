package parser

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

// Structs for Version 2 & 3
type SwiftPinV2 struct {
	Identity string `json:"identity"`
	Location string `json:"location"`
	State    struct {
		Version string `json:"version"`
	} `json:"state"`
}

type SwiftResolvedV2 struct {
	Pins    []SwiftPinV2 `json:"pins"`
	Version int          `json:"version"`
}

// Structs for Version 1
type SwiftPinV1 struct {
	Package       string `json:"package"`
	RepositoryURL string `json:"repositoryURL"`
	State         struct {
		Version string `json:"version"`
	} `json:"state"`
}

type SwiftResolvedV1 struct {
	Object struct {
		Pins []SwiftPinV1 `json:"pins"`
	} `json:"object"`
	Version int `json:"version"`
}

type SwiftVersionCheck struct {
	Version int `json:"version"`
}

func ParseSwiftResolved(path string) ([]models.Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ver SwiftVersionCheck
	if err := json.Unmarshal(data, &ver); err != nil {
		return nil, err
	}

	var deps []models.Dependency

	if ver.Version == 1 {
		var resolved SwiftResolvedV1
		if err := json.Unmarshal(data, &resolved); err != nil {
			return nil, err
		}
		for _, pin := range resolved.Object.Pins {
			name := pin.Package
			location := pin.RepositoryURL
			if name == "" {
				parts := strings.Split(strings.TrimSuffix(location, ".git"), "/")
				if len(parts) > 0 {
					name = parts[len(parts)-1]
				}
			}

			docURL := cleanGitURL(location)

			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   pin.State.Version,
				Ecosystem: "swift",
				Type:      "prod",
				DocURL:    docURL,
			})
		}
	} else {
		var resolved SwiftResolvedV2
		if err := json.Unmarshal(data, &resolved); err != nil {
			return nil, err
		}
		for _, pin := range resolved.Pins {
			name := pin.Identity
			location := pin.Location
			if name == "" {
				parts := strings.Split(strings.TrimSuffix(location, ".git"), "/")
				if len(parts) > 0 {
					name = parts[len(parts)-1]
				}
			}

			docURL := cleanGitURL(location)

			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   pin.State.Version,
				Ecosystem: "swift",
				Type:      "prod",
				DocURL:    docURL,
			})
		}
	}

	return deps, nil
}

func cleanGitURL(gitURL string) string {
	url := strings.TrimSpace(gitURL)
	url = strings.TrimSuffix(url, ".git")

	// If it is SSH format (git@github.com:owner/repo), convert to HTTPS
	if strings.HasPrefix(url, "git@") {
		url = strings.Replace(url, ":", "/", 1)
		url = strings.Replace(url, "git@", "https://", 1)
	}

	return url
}
