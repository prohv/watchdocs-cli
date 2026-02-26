package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "watchdocs",
	Short: "WatchDocs CLI - Find docs instantly",
}

func Execute() error {
	return rootCmd.Execute()
}