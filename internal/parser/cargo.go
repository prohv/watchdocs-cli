package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParseCargoToml(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	var deps []models.Dependency
	lines := strings.Split(content, "\n")

	reSimple := regexp.MustCompile(`^(\S+)\s*=\s*"([^"]+)"`)
	reTable := regexp.MustCompile(`^([\w\-.]+)\s*=\s*\{`)
	reVersion := regexp.MustCompile(`version\s*=\s*"([^"]+)"`)
	rePackage := regexp.MustCompile(`package\s*=\s*"([^"]+)"`)
	reDevDeps := regexp.MustCompile(`dev[-_]dependencies`)
	reTestDeps := regexp.MustCompile(`test[-_]dependencies`)

	var currentType string // "prod" | "dev" | "test" | ""

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "[") {
			if strings.Contains(trimmed, "dependencies") {
				if reDevDeps.MatchString(trimmed) {
					currentType = "dev"
				} else if reTestDeps.MatchString(trimmed) {
					currentType = "test"
				} else {
					currentType = "prod"
				}
			} else {
				currentType = ""
			}
			continue
		}

		if currentType == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		if m := reSimple.FindStringSubmatch(trimmed); m != nil {
			deps = append(deps, models.Dependency{
				Name:      m[1],
				Version:   m[2],
				Ecosystem: "cargo",
				Type:      currentType,
			})
			continue
		}

		if m := reTable.FindStringSubmatch(trimmed); m != nil {
			name := m[1]
			if pkg := rePackage.FindStringSubmatch(trimmed); pkg != nil {
				name = pkg[1]
			}
			version := ""
			if ver := reVersion.FindStringSubmatch(trimmed); ver != nil {
				version = ver[1]
			}
			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   version,
				Ecosystem: "cargo",
				Type:      currentType,
			})
		}
	}

	return deps, nil
}
