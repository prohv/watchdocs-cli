package parser

import (
	"encoding/json"
	"os"
)

type PackageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func ParseNPM(path string) ([]Dependency, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: version})
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, Dependency{Name: name, Version: version})
	}

	return deps, nil
}
