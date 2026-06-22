package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

var reDepLine = regexp.MustCompile(`^\s+([a-zA-Z0-9_.-]+):(.*)`)

func ParsePubspecYaml(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	var deps []models.Dependency
	seen := map[string]bool{}
	lines := strings.Split(content, "\n")
	currentType := "" // "prod" | "dev" | ""

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// top-level block detection (no leading whitespace)
		if strings.HasPrefix(line, "dependencies:") {
			currentType = "prod"
			continue
		}
		if strings.HasPrefix(line, "dev_dependencies:") {
			currentType = "dev"
			continue
		}
		// any other top-level key resets
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' {
			currentType = ""
			continue
		}

		if currentType == "" {
			continue
		}

		// indented dep line
		depMatch := reDepLine.FindStringSubmatch(line)
		if depMatch == nil {
			continue
		}

		name := depMatch[1]
		version := strings.TrimSpace(depMatch[2])

		// skip SDK / flutter pseudo-deps
		if name == "flutter" || name == "sdk" {
			continue
		}

		// skip if next line is `sdk:`
		if version == "" || version == "any" || strings.HasPrefix(version, "{") {
			version = ""
			if i+1 < len(lines) && strings.HasPrefix(strings.TrimSpace(lines[i+1]), "sdk:") {
				continue
			}
		} else {
			version = strings.Trim(version, `'"`)
		}

		if !seen[name] {
			seen[name] = true
			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   version,
				Ecosystem: "pub",
				Type:      currentType,
			})
		}
	}

	return deps, nil
}
