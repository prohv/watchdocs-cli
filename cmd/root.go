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

func init() {
	rootCmd.SetHelpTemplate(`WatchDocs CLI - Find docs instantly

Usage:
  watchdocs <command> [flags]

Commands:
  scan [flags]                 Scan project recursively for dependency doc URLs
    -p, --path <dir>           Target directory (default: cwd)
    -e, --ecosystem <list>     Filter to specific ecosystems, e.g. npm,go
    -s, --slim                 Return only name + docUrl (saves tokens)

  lookup <package> [flags]     Lookup a single package without a manifest
    -e, --ecosystem <eco>      Required: npm | go | pip | cargo | pub | maven
    -s, --slim                 Return only name + docUrl (saves tokens)

Flags:
  -h, --help                   Show this help

Supported ecosystems:  npm · go · pip · cargo · pub · maven
Supported manifests:   package.json · go.mod · requirements.txt · pyproject.toml · uv.lock · Cargo.toml · pom.xml · pubspec.yaml
`)
}