package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/prohv/watchdocs-cli/internal/models"
	"github.com/prohv/watchdocs-cli/internal/parser"
	"github.com/prohv/watchdocs-cli/internal/resolver"
	"github.com/prohv/watchdocs-cli/internal/scanner"
	"github.com/spf13/cobra"
)

type ScanResult struct {
	Scanned []string      `json:"scanned"`
	Total   int           `json:"total"`
	Results []DepResult   `json:"results"`
}

type DepResult struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	Ecosystem string `json:"ecosystem"`
	Type      string `json:"type"`
	DocURL    string `json:"docUrl,omitempty"`
	Status    string `json:"status"` // "resolved" | "not_found"
}

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project for dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		cwd, _ := os.Getwd()
		found, err := scanner.Scan(cwd)
		if err != nil {
			printError("scan_failed", err.Error())
			return
		}

		if len(found) == 0 {
			printError("no_manifests", "no manifest files found in project root")
			return
		}

		var scanned []string
		var allDeps []models.Dependency

		for _, m := range found {
			scanned = append(scanned, m.Type)

			content, err := os.ReadFile(m.Path)
			if err != nil {
				printError("read_failed", fmt.Sprintf("could not read %s: %v", m.Type, err))
				continue
			}

			if m.Type == "package.json" {
				deps, err := parser.ParseNPM(m.Path)
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "go.mod" {
				deps, err := parser.ParseGoMod(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "Cargo.toml" {
				deps, err := parser.ParseCargoToml(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "requirements.txt" {
				deps, err := parser.ParseRequirementsTxt(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "pyproject.toml" {
				deps, err := parser.ParsePyProjectToml(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "pubspec.yaml" {
				deps, err := parser.ParsePubspecYaml(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "pom.xml" {
				deps, err := parser.ParsePomXml(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "uv.lock" {
				deps, err := parser.ParseUvLock(string(content))
				if err != nil {
					printError("parse_failed", fmt.Sprintf("could not parse %s: %v", m.Type, err))
					continue
				}
				allDeps = append(allDeps, deps...)
			}
		}

		if len(allDeps) == 0 {
			printError("no_deps", "no dependencies found in detected manifests")
			return
		}

		var results []DepResult

		for _, dep := range allDeps {
			var result *models.DocResult

			switch dep.Ecosystem {
			case "npm":
				result, err = resolver.OnlineNpmResolver(dep)
			case "go":
				result, err = resolver.OnlineGoResolver(dep)
			case "pip":
				result, err = resolver.OnlinePipResolver(dep)
			case "cargo":
				result, err = resolver.OnlineCargoResolver(dep)
			case "pub":
				result, err = resolver.OnlinePubResolver(dep)
			case "maven":
				result, err = resolver.OnlineMavenResolver(dep)
			}

			if err != nil || result == nil || result.DocURL == "" {
				results = append(results, DepResult{
					Name:      dep.Name,
					Version:   dep.Version,
					Ecosystem: dep.Ecosystem,
					Type:      dep.Type,
					Status:    "not_found",
				})
				continue
			}

			results = append(results, DepResult{
				Name:      result.Name,
				Version:   result.Version,
				Ecosystem: result.Ecosystem,
				Type:      result.Type,
				DocURL:    result.DocURL,
				Status:    "resolved",
			})
		}

		output := ScanResult{
			Scanned: scanned,
			Total:   len(results),
			Results: results,
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(output)
	},
}

func printError(code string, msg string) {
	type errOut struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(errOut{Error: code, Message: msg})
}