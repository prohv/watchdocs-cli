package parser

import (
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParseGoMod(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	lines := strings.Split(content, "\n")
	var deps []models.Dependency
	inRequireBlock := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if trimmed == "require (" {
			inRequireBlock = true
			continue
		}

		if trimmed == ")" {
			inRequireBlock = false
			continue
		}

		if match := matchSingleRequire(trimmed); match != nil {
			deps = append(deps, *match)
			continue
		}

		if inRequireBlock && !strings.HasPrefix(trimmed, "//") {
			if match := matchBlockRequire(trimmed); match != nil {
				deps = append(deps, *match)
			}
		}
	}

	return deps, nil
}

func matchSingleRequire(line string) *models.Dependency {
	if !strings.HasPrefix(line, "require ") {
		return nil
	}
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return nil
	}
	return &models.Dependency{Name: parts[1], Version: parts[2], Ecosystem: "go", Type: "prod"}
}

func matchBlockRequire(line string) *models.Dependency {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil
	}
	return &models.Dependency{Name: parts[0], Version: parts[1], Ecosystem: "go", Type: "prod"}
}
