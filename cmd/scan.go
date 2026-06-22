package cmd

import (
	"fmt"
	"os"

	"github.com/prohv/watchdocs-cli/internal/models"
	"github.com/prohv/watchdocs-cli/internal/parser"
	"github.com/prohv/watchdocs-cli/internal/resolver"
	"github.com/prohv/watchdocs-cli/internal/scanner"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project for dependencies",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning project...")
		cwd, _ := os.Getwd()
		found, err := scanner.Scan(cwd)

		if err != nil {
			fmt.Println("Error scanning:", err)
			return
		}

		if len(found) == 0 {
			fmt.Println("No manifest files found!")
			return
		}

		var allDeps []models.Dependency
		for _, m := range found {
			fmt.Printf("Found: %s\n", m.Type)

			content, err := os.ReadFile(m.Path)
			if err != nil {
				fmt.Printf("Error reading %s: %v\n", m.Type, err)
				continue
			}

			if m.Type == "package.json" {
				deps, err := parser.ParseNPM(m.Path)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", m.Type, err)
					continue
				}
				allDeps = append(allDeps, deps...)
			} else if m.Type == "go.mod" {
				deps, err := parser.ParseGoMod(string(content))
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", m.Type, err)
					continue
				}
				allDeps = append(allDeps, deps...)
			}
		}

		if len(allDeps) == 0 {
			fmt.Println("No dependencies found.")
			return
		}

		fmt.Println("\n--- Resolving docs ---")
		for _, dep := range allDeps {
			var result *models.DocResult

			switch dep.Ecosystem {
			case "npm":
				result, err = resolver.OnlineNpmResolver(dep)
			case "go":
				result, err = resolver.OnlineGoResolver(dep)
			}

			if err != nil || result == nil || result.DocURL == "" {
				fmt.Printf("%-40s -> (not found)\n", dep.Name)
				continue
			}
			fmt.Printf("%-40s -> %s\n", result.Name, result.DocURL)
		}
	},
}