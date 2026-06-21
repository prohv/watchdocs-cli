package parser

import (
	"encoding/json"
	"os"

	"github.com/prohv/watchdocs-cli/internal/models"
)

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParseNPM(path string) ([]models.Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []models.Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, models.Dependency{Name: name, Version: version, Ecosystem: "npm", Type: "prod"})
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, models.Dependency{Name: name, Version: version, Ecosystem: "npm", Type: "dev"})
	}

	return deps, nil
}
