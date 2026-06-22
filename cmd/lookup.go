package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	lookupCmd.Flags().StringP("ecosystem", "e", "", "Ecosystem to resolve against (npm, go, pip, cargo, pub, maven)")
	lookupCmd.Flags().BoolP("slim", "s", false, "Return only name and docUrl (saves tokens)")
	lookupCmd.MarkFlagRequired("ecosystem")
	rootCmd.AddCommand(lookupCmd)
}

var lookupCmd = &cobra.Command{
	Use:   "lookup <package>",
	Short: "Lookup doc URL for a single package without a manifest",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.TrimSpace(args[0])
		ecosystem, _ := cmd.Flags().GetString("ecosystem")
		ecosystem = strings.TrimSpace(strings.ToLower(ecosystem))

		validEcosystems := map[string]bool{
			"npm": true, "go": true, "pip": true,
			"cargo": true, "pub": true, "maven": true,
		}

		if !validEcosystems[ecosystem] {
			printError("invalid_ecosystem", "must be one of: npm, go, pip, cargo, pub, maven")
			return
		}

		slim, _ := cmd.Flags().GetBool("slim")

		dep := models.Dependency{
			Name:      name,
			Ecosystem: ecosystem,
		}

		result := resolveDoc(dep)

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")

		if slim {
			type slimResult struct {
				Name   string `json:"name"`
				DocURL string `json:"docUrl,omitempty"`
			}
			enc.Encode(slimResult{Name: result.Name, DocURL: result.DocURL})
			return
		}

		enc.Encode(result)
	},
}
