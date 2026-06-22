package parser

import (
	"regexp"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParsePomXml(content string) ([]models.Dependency, error) {
	if strings.TrimSpace(content) == "" {
		return []models.Dependency{}, nil
	}

	reDep := regexp.MustCompile(`(?i)<dependency>([\s\S]*?)<\/dependency>`)
	reGroupID := regexp.MustCompile(`(?i)<groupId>\s*(.*?)\s*<\/groupId>`)
	reArtifactID := regexp.MustCompile(`(?i)<artifactId>\s*(.*?)\s*<\/artifactId>`)
	reVersion := regexp.MustCompile(`(?i)<version>\s*(.*?)\s*<\/version>`)
	reScope := regexp.MustCompile(`(?i)<scope>\s*(.*?)\s*<\/scope>`)

	var deps []models.Dependency

	for _, match := range reDep.FindAllStringSubmatch(content, -1) {
		block := match[1]

		groupMatch := reGroupID.FindStringSubmatch(block)
		artifactMatch := reArtifactID.FindStringSubmatch(block)

		if groupMatch == nil || artifactMatch == nil {
			continue
		}

		name := groupMatch[1] + ":" + artifactMatch[1]

		version := ""
		if v := reVersion.FindStringSubmatch(block); v != nil {
			version = v[1]
			if strings.HasPrefix(version, "${") {
				version = "" // unresolved property placeholder
			}
		}

		depType := "prod"
		if s := reScope.FindStringSubmatch(block); s != nil {
			switch strings.ToLower(s[1]) {
			case "test":
				depType = "test"
			case "provided":
				depType = "dev"
			}
		}

		deps = append(deps, models.Dependency{
			Name:      name,
			Version:   version,
			Ecosystem: "maven",
			Type:      depType,
		})
	}

	return deps, nil
}
