package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
	"github.com/prohv/watchdocs-cli/internal/parser"
	"github.com/prohv/watchdocs-cli/internal/resolver"
	"github.com/prohv/watchdocs-cli/internal/scanner"
	"github.com/spf13/cobra"
)

type ScanResult struct {
	Scanned []string    `json:"scanned"`
	Total   int         `json:"total"`
	Results []DepResult `json:"results"`
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
	scanCmd.Flags().StringP("path", "p", "", "Path to project directory (defaults to cwd)")
	scanCmd.Flags().StringP("ecosystem", "e", "", "Filter by ecosystem(s), comma-separated (e.g. npm,go)")
	scanCmd.Flags().BoolP("slim", "s", false, "Return only name and docUrl per result (saves tokens)")
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project for dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		pathFlag, _ := cmd.Flags().GetString("path")
		ecoFlag, _ := cmd.Flags().GetString("ecosystem")
		slim, _ := cmd.Flags().GetBool("slim")

		root := pathFlag
		if root == "" {
			root, _ = os.Getwd()
		}

		// build ecosystem filter set
		ecoFilter := map[string]bool{}
		if ecoFlag != "" {
			for _, e := range strings.Split(ecoFlag, ",") {
				ecoFilter[strings.TrimSpace(e)] = true
			}
		}

		found, err := scanner.Scan(root)
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

		// deduplicate by ecosystem:name
		seen := map[string]bool{}
		var uniqueDeps []models.Dependency
		for _, dep := range allDeps {
			key := dep.Ecosystem + ":" + dep.Name
			if !seen[key] {
				seen[key] = true
				uniqueDeps = append(uniqueDeps, dep)
			}
		}
		allDeps = uniqueDeps

		// apply ecosystem filter if set
		if len(ecoFilter) > 0 {
			var filtered []models.Dependency
			for _, dep := range allDeps {
				if ecoFilter[dep.Ecosystem] {
					filtered = append(filtered, dep)
				}
			}
			allDeps = filtered
		}

		if len(allDeps) == 0 {
			printError("no_deps", "no dependencies found for the specified ecosystem(s)")
			return
		}

		var results []DepResult

		for _, dep := range allDeps {
			result := resolveDoc(dep)
			results = append(results, result)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		if slim {
			type slimResult struct {
				Name   string `json:"name"`
				DocURL string `json:"docUrl,omitempty"`
			}
			type slimOutput struct {
				Scanned []string    `json:"scanned"`
				Total   int         `json:"total"`
				Results []slimResult `json:"results"`
			}
			var slim []slimResult
			for _, r := range results {
				slim = append(slim, slimResult{Name: r.Name, DocURL: r.DocURL})
			}
			enc.Encode(slimOutput{Scanned: scanned, Total: len(slim), Results: slim})
			return
		}

		enc.Encode(ScanResult{Scanned: scanned, Total: len(results), Results: results})
	},
}

func resolveDoc(dep models.Dependency) DepResult {
	var result *models.DocResult
	var err error

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
		return DepResult{
			Name:      dep.Name,
			Version:   dep.Version,
			Ecosystem: dep.Ecosystem,
			Type:      dep.Type,
			Status:    "not_found",
		}
	}

	return DepResult{
		Name:      result.Name,
		Version:   result.Version,
		Ecosystem: result.Ecosystem,
		Type:      result.Type,
		DocURL:    result.DocURL,
		Status:    "resolved",
	}
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