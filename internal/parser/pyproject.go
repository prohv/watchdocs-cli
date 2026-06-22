package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParsePyProjectToml(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	var result []models.Dependency
	lines := strings.Split(content, "\n")
	section := ""

	reSection := regexp.MustCompile(`^\[(.+)\]$`)
	reKV := regexp.MustCompile(`^([\w\-.]+)\s*=\s*(.+)$`)
	reVersion := regexp.MustCompile(`version\s*=\s*["']([^"']+)["']`)
	reGroupDep := regexp.MustCompile(`^(\w+)\s*=\s*\[([^\]]*)\]`)

	for _, raw := range lines {
		line := strings.TrimSpace(raw)

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		if sec := reSection.FindStringSubmatch(line); sec != nil {
			section = sec[1]
			continue
		}

		// Poetry deps
		if section == "tool.poetry.dependencies" || section == "tool.poetry.dev-dependencies" {
			m := reKV.FindStringSubmatch(line)
			if m == nil || m[1] == "python" {
				continue
			}
			version := strings.Trim(m[2], `"'`)
			if v := reVersion.FindStringSubmatch(version); v != nil {
				version = v[1]
			}
			depType := "prod"
			if section == "tool.poetry.dev-dependencies" {
				depType = "dev"
			}
			result = append(result, models.Dependency{Name: m[1], Version: version, Ecosystem: "pip", Type: depType})
		}

		// PEP 621 project.dependencies inline array
		if section == "project" && strings.HasPrefix(line, "dependencies") {
			re := regexp.MustCompile(`dependencies\s*=\s*\[([\s\S]*?)\]`)
			if m := re.FindStringSubmatch(content); m != nil {
				for _, item := range strings.Split(m[1], ",") {
					item = strings.Trim(strings.TrimSpace(item), `"'`)
					if item == "" {
						continue
					}
					name, version := splitPyDep(item)
					result = append(result, models.Dependency{Name: name, Version: version, Ecosystem: "pip", Type: "prod"})
				}
			}
		}

		// Optional dependency groups
		if strings.HasPrefix(section, "project.optional-dependencies") {
			parts := strings.Split(section, ".")
			groupName := strings.ToLower(parts[len(parts)-1])
			depType := "dev"
			if strings.Contains(groupName, "test") {
				depType = "test"
			}
			if m := reKV.FindStringSubmatch(line); m != nil {
				result = append(result, models.Dependency{Name: m[1], Version: m[2], Ecosystem: "pip", Type: depType})
			}
		}

		// dependency-groups (PEP 735)
		if section == "dependency-groups" {
			if m := reGroupDep.FindStringSubmatch(line); m != nil {
				groupName := strings.ToLower(m[1])
				depType := "dev"
				if strings.Contains(groupName, "test") {
					depType = "test"
				}
				for _, item := range strings.Split(m[2], ",") {
					item = strings.Trim(strings.TrimSpace(item), `"'`)
					if item == "" {
						continue
					}
					name, version := splitPyDep(item)
					result = append(result, models.Dependency{Name: name, Version: version, Ecosystem: "pip", Type: depType})
				}
			}
		}
	}

	return result, nil
}

func splitPyDep(dep string) (string, string) {
	re := regexp.MustCompile(`^([\w\-.]+)(.*)$`)
	m := re.FindStringSubmatch(dep)
	if m == nil {
		return dep, "unknown"
	}
	version := strings.TrimSpace(m[2])
	if version == "" {
		return m[1], "unknown"
	}
	return m[1], version
}
