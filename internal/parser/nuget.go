package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParseCSProj(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	// Match <PackageReference ... /> or <PackageVersion ... />
	reRef := regexp.MustCompile(`(?i)<(PackageReference|PackageVersion)\s+([\s\S]*?)(\/>|>)`)
	reIncludeAttr := regexp.MustCompile(`(?i)Include\s*=\s*["']\s*(.*?)\s*["']`)
	reVersionAttr := regexp.MustCompile(`(?i)Version\s*=\s*["']\s*(.*?)\s*["']`)

	var deps []models.Dependency

	matches := reRef.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		tagBody := match[2]

		includeMatch := reIncludeAttr.FindStringSubmatch(tagBody)
		if includeMatch == nil {
			continue
		}
		name := includeMatch[1]

		version := ""
		versionMatch := reVersionAttr.FindStringSubmatch(tagBody)
		if versionMatch != nil {
			version = versionMatch[1]
		}

		if name == "" {
			continue
		}

		depType := "prod"
		nameLower := strings.ToLower(name)
		if strings.Contains(nameLower, "test") || strings.Contains(nameLower, "xunit") || strings.Contains(nameLower, "nunit") {
			depType = "dev"
		}

		deps = append(deps, models.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "nuget",
			Type:      depType,
		})
	}

	return deps, nil
}

func ParsePackagesConfig(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	rePackage := regexp.MustCompile(`(?i)<package\s+([\s\S]*?)\/>`)
	reIdAttr := regexp.MustCompile(`(?i)id\s*=\s*["']\s*(.*?)\s*["']`)
	reVersionAttr := regexp.MustCompile(`(?i)version\s*=\s*["']\s*(.*?)\s*["']`)

	var deps []models.Dependency

	for _, match := range rePackage.FindAllStringSubmatch(content, -1) {
		tagBody := match[1]
		idMatch := reIdAttr.FindStringSubmatch(tagBody)
		if idMatch == nil {
			continue
		}
		name := idMatch[1]

		version := ""
		versionMatch := reVersionAttr.FindStringSubmatch(tagBody)
		if versionMatch != nil {
			version = versionMatch[1]
		}

		if name == "" {
			continue
		}

		depType := "prod"
		nameLower := strings.ToLower(name)
		if strings.Contains(nameLower, "test") || strings.Contains(nameLower, "xunit") || strings.Contains(nameLower, "nunit") {
			depType = "dev"
		}

		deps = append(deps, models.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "nuget",
			Type:      depType,
		})
	}

	return deps, nil
}
