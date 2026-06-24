package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/prohv/watchdocs-cli/internal/cache"
	"github.com/prohv/watchdocs-cli/internal/models"
	"github.com/spf13/cobra"
)

func init() {
	lookupCmd.Flags().StringP("ecosystem", "e", "", "Ecosystem to resolve against (npm, go, pip, cargo, pub, maven)")
	lookupCmd.Flags().BoolP("slim", "s", false, "Return only name and docUrl (saves tokens)")
	lookupCmd.Flags().Bool("clear-cache", false, "Clear the local cache before looking up")
	lookupCmd.Flags().Bool("no-cache", false, "Disable reading/writing to the local cache")
	lookupCmd.Flags().StringP("format", "f", "json", "Output format (json | list)")
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
		clearCache, _ := cmd.Flags().GetBool("clear-cache")
		noCache, _ := cmd.Flags().GetBool("no-cache")
		formatFlag, _ := cmd.Flags().GetString("format")
		formatFlag = strings.ToLower(strings.TrimSpace(formatFlag))

		validEcosystems := map[string]bool{
			"npm": true, "go": true, "pip": true,
			"cargo": true, "pub": true, "maven": true,
		}

		if !validEcosystems[ecosystem] {
			printError("invalid_ecosystem", "must be one of: npm, go, pip, cargo, pub, maven")
			return
		}

		slim, _ := cmd.Flags().GetBool("slim")

		// Initialize Cache
		var c *cache.Cache
		if !noCache {
			var err error
			c, err = cache.NewCache()
			if err == nil && clearCache {
				_ = c.Clear()
			}
		}

		dep := models.Dependency{
			Name:      name,
			Ecosystem: ecosystem,
		}

		var result DepResult
		cacheHit := false

		if !noCache && c != nil {
			if cachedURL, found := c.Get(dep.Ecosystem, dep.Name); found {
				result = DepResult{
					Name:      dep.Name,
					Ecosystem: dep.Ecosystem,
					DocURL:    cachedURL,
					Status:    "resolved",
				}
				cacheHit = true
			}
		}

		if !cacheHit {
			result = resolveDoc(dep)
			if !noCache && c != nil && result.Status == "resolved" && result.DocURL != "" {
				cacheURL := result.DocURL
				if dep.Ecosystem == "go" {
					cacheURL = strings.SplitN(result.DocURL, "@", 2)[0]
				}
				c.Set(dep.Ecosystem, dep.Name, cacheURL)
				_ = c.Save()
			}
		}

		if formatFlag == "list" || formatFlag == "l" {
			printLookupList(result, slim)
			return
		}

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

func printLookupList(result DepResult, slim bool) {
	docLink := "Not resolved"
	if result.DocURL != "" {
		docLink = fmt.Sprintf("[docs](%s)", result.DocURL)
	}
	if slim {
		fmt.Printf("- **%s** - %s\n", result.Name, docLink)
	} else {
		fmt.Printf("- **%s** (%s, %s) - %s\n", result.Name, result.Ecosystem, result.Status, docLink)
	}
}
