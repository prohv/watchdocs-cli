package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/prohv/watchdocs-cli/internal/scanner"
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

		for _, m := range found {
			fmt.Printf("Found: %s\n", m.Type)
		}
	},
}