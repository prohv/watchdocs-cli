package scanner

import (
	"os"
	"path/filepath"
)

type DetectedManifest struct {
	Type string
	Path string
}

func Scan(root string) ([]DetectedManifest, error) {
	manifests := []string{
		"package.json", "go.mod", "requirements.txt",
		"pyproject.toml", "Cargo.toml", "pom.xml",
	}

	var found []DetectedManifest

	for _, m := range manifests {
		path := filepath.Join(root, m)
		if _, err := os.Stat(path); err == nil {
			found = append(found, DetectedManifest{
				Type: m,
				Path: path,
			})
		}
	}
	return found, nil
}