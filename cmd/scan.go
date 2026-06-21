package cmd

import (
	"fmt"
	"os"

	"github.com/prohv/watchdocs-cli/internal/parser"
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
			fmt.Println("No dependencies found.")
			return
		}

		fmt.Println("\n--- Dependencies Found ---")
		for _, dep := range allDeps {
			fmt.Println(" -", dep)
		}
	},
}