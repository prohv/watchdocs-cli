package scanner

import (
	"io/fs"
	"path/filepath"
)

type DetectedManifest struct {
	Type string
	Path string
}

var manifests = map[string]bool{
	"package.json":  true,
	"go.mod":        true,
	"requirements.txt": true,
	"pyproject.toml": true,
	"Cargo.toml":    true,
	"pom.xml":       true,
	"uv.lock":       true,
	"pubspec.yaml":  true,
}

// heavy directories to skip — equivalent to VS Code's exclude globs
var skipDirs = map[string]bool{
	"node_modules": true,
	".git":         true,
	"vendor":       true,
	"__pycache__":  true,
	".venv":        true,
	"venv":         true,
	"env":          true,
	"target":       true,
	"dist":         true,
	"build":        true,
	".cache":       true,
	".idea":        true,
	".vscode":      true,
	"out":          true,
	"bin":          true,
}

func Scan(root string) ([]DetectedManifest, error) {
	var found []DetectedManifest

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}

		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		if manifests[d.Name()] {
			found = append(found, DetectedManifest{
				Type: d.Name(),
				Path: path,
			})
		}

		return nil
	})

	return found, err
}