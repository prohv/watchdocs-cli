package parser

import (
	"bufio"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
)

func ParseGoMod(path string) ([]models.Dependency, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var deps []models.Dependency
	scanner := bufio.NewScanner(file)
	inRequire := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if !inRequire && (strings.HasPrefix(line, "module") || strings.HasPrefix(line, "go ")) {
			continue
		}

		if line == "require (" {
			inRequire = true
			continue
		}

		if line == ")" {
			inRequire = false
			continue
		}

		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps = append(deps, models.Dependency{Name: parts[1], Version: parts[2], Ecosystem: "go"})
			}
			continue
		}

		if inRequire {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				deps = append(deps, models.Dependency{Name: parts[0], Version: parts[1], Ecosystem: "go"})
			}
		}
	}

	return deps, scanner.Err()
}
