package cli

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/spf13/cobra"
)

// ProviderPathCompletion provides shell completion for paths that may have a provider prefix (e.g., "local:/tmp").
// It handles the case where Bash splits words at colons by reconstructing the full path from args.
func ProviderPathCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		currentToComplete := toComplete
		prefixToStrip := ""

		// reconstruct currentToComplete if prefixToStrip is set
		if prefixToStrip != "" {
			currentToComplete = prefixToStrip + toComplete
		}

		// 1. Determine if we are completing a provider or a path
		uri, err := fs.ParseURI(currentToComplete)
		if err != nil {
			// If parsing fails, it might be an incomplete provider name
			if !strings.Contains(currentToComplete, ":") && !strings.HasPrefix(currentToComplete, "/") {
				names, err := container.ProviderRegistry().RegisteredNames()
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}

				var results []string
				for _, name := range names {
					if strings.HasPrefix(name, currentToComplete) {
						results = append(results, name+":")
					}
				}
				return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		provider := uri.Provider
		path := uri.Path
		found := strings.Contains(currentToComplete, ":")

		// If path is empty (just "provider:"), suggest "/"
		if path == "" {
			res := "/"
			if found && prefixToStrip == "" {
				res = provider + ":" + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		// Find parent directory and prefix
		lastSlash := strings.LastIndex(path, "/")
		if lastSlash == -1 {
			// No slash after the colon (e.g. "local:tm") or path is just "tm"
			res := "/"
			if found && prefixToStrip == "" {
				res = provider + ":" + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		listDir := path[:lastSlash+1]
		searchPrefix := path[lastSlash+1:]

		// Trim trailing slash for validation (except for root "/")
		cleanListDir := listDir
		if len(cleanListDir) > 1 && strings.HasSuffix(cleanListDir, "/") {
			cleanListDir = strings.TrimSuffix(cleanListDir, "/")
		}

		// Full path for manager to resolve correctly
		managerListDir := provider + ":" + cleanListDir

		items, err := container.FS().List(cmd.Context(), managerListDir, fs.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, item := range items {
			if strings.HasPrefix(item.Name, searchPrefix) {
				res := listDir + item.Name
				if item.Type == fs.TypeFolder {
					res += "/"
				}

				// If we reconstructed the prefix from 'args' (meaning the shell split it),
				// we return the path without the prefix so the shell matches it correctly.
				// If the prefix was in 'toComplete', we return the full path.
				if found && prefixToStrip == "" {
					results = append(results, provider+":"+res)
				} else {
					results = append(results, res)
				}
			}
		}

		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
