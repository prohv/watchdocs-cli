package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParseUvLock(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	var deps []models.Dependency
	lines := strings.Split(content, "\n")

	inPackage := false
	name := ""
	version := "unknown"

	reName := regexp.MustCompile(`^name\s*=\s*"([^"]+)"`)
	reVersion := regexp.MustCompile(`^version\s*=\s*"([^"]+)"`)

	flush := func() {
		if inPackage && name != "" {
			deps = append(deps, models.Dependency{
				Name:      name,
				Version:   version,
				Ecosystem: "pip",
				Type:      classifyPyDep(strings.ToLower(name)),
			})
		}
		inPackage = false
		name = ""
		version = "unknown"
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if line == "[[package]]" {
			flush()
			inPackage = true
			continue
		}

		if strings.HasPrefix(line, "[") && !strings.HasPrefix(line, "[[") {
			flush()
			continue
		}

		if inPackage {
			if m := reName.FindStringSubmatch(line); m != nil {
				name = m[1]
				continue
			}
			if m := reVersion.FindStringSubmatch(line); m != nil {
				version = m[1]
				continue
			}
		}
	}

	flush()

	return deps, nil
}
