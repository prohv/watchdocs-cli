package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/prohv/watchdocs-cli/internal/scanner"
	"github.com/prohv/watchdocs-cli/internal/resolver"
	"github.com/prohv/watchdocs-cli/internal/parser"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project for dependencies and find docs",
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
		var allDeps []string
		for _, m := range found {
			fmt.Printf("Found: %s\n", m.Type)
			if m.Type == "package.json" {
				deps, err := parser.ParseNPM(m.Path)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", m.Type, err)
					continue
				}
				for _, d := range deps {
					allDeps = append(allDeps, d.Name)
				}
			} else if m.Type == "go.mod" {
				deps, err := parser.ParseGoMod(m.Path)
				if err != nil {
					fmt.Printf("Error parsing %s: %v\n", m.Type, err)
					continue
				}
				for _, d := range deps {
					allDeps = append(allDeps, d.Name)
				}
			}
		}

		if len(allDeps) == 0 {
			fmt.Println("No dependencies found to resolve.")
			return
		}

		fmt.Printf("Asking Gemini for docs for %d dependencies...\n", len(allDeps))
		results, err := resolver.ResolveDocs(allDeps)
		if err != nil {
			fmt.Printf("API Error: %v\n", err)
			return
		}

		fmt.Println("\n--- Official Documentation ---")
		for _, res := range results {
			fmt.Printf("%-20s -> %s\n", res.Name, res.DocURL)
		}
	},
}