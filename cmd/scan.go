package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(scanCmd)
}

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan project for dependencies and find docs",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Scanning project...")
		// Logic coming soon!
	},
}