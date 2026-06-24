package parser

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

type ComposerJSON struct {
	Require    map[string]string `json:"require"`
	RequireDev map[string]string `json:"require-dev"`
}

func ParseComposer(path string) ([]models.Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg ComposerJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []models.Dependency
	for name, version := range pkg.Require {
		if name == "php" || strings.HasPrefix(name, "ext-") {
			continue
		}
		deps = append(deps, models.Dependency{Name: name, Version: version, Ecosystem: "composer", Type: "prod"})
	}
	for name, version := range pkg.RequireDev {
		if name == "php" || strings.HasPrefix(name, "ext-") {
			continue
		}
		deps = append(deps, models.Dependency{Name: name, Version: version, Ecosystem: "composer", Type: "dev"})
	}

	return deps, nil
}
