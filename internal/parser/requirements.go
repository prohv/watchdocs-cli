package parser

import (
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

var operators = []string{"==", ">=", "<=", "!=", "~=", ">", "<"}

func ParseRequirementsTxt(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	var deps []models.Dependency

	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		name := trimmed
		version := "unknown"

		for _, op := range operators {
			if strings.Contains(trimmed, op) {
				parts := strings.SplitN(trimmed, op, 2)
				name = strings.TrimSpace(parts[0])
				version = strings.TrimSpace(parts[1])
				break
			}
		}

		// strip extras e.g. requests[security]
		name = strings.TrimSpace(strings.Split(name, "[")[0])

		depType := classifyPyDep(strings.ToLower(name))

		deps = append(deps, models.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "pip",
			Type:      depType,
		})
	}

	return deps, nil
}

func classifyPyDep(name string) string {
	testKeywords := []string{"pytest", "tox", "coverage", "unittest"}
	devKeywords := []string{"black", "flake8", "mypy", "pylint", "isort", "ruff", "pre-commit"}

	for _, kw := range testKeywords {
		if strings.Contains(name, kw) {
			return "test"
		}
	}
	for _, kw := range devKeywords {
		if strings.Contains(name, kw) {
			return "dev"
		}
	}
	return "prod"
}
