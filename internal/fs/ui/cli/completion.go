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

		// Attempt to reconstruct the current argument if Bash split it on colons.
		// We look at the tail of 'args' to see if it ends with a provider prefix.
		if !strings.Contains(toComplete, ":") && len(args) > 0 {
			if args[len(args)-1] == ":" && len(args) >= 2 {
				prefixToStrip = args[len(args)-2] + ":"
				currentToComplete = prefixToStrip + toComplete
			} else if strings.HasSuffix(args[len(args)-1], ":") {
				prefixToStrip = args[len(args)-1]
				currentToComplete = prefixToStrip + toComplete
			}
		}

		// 1. Determine if we are completing a provider or a path
		provider, path, found := fs.SplitProviderPath(currentToComplete)

		// If no colon and doesn't start with a slash, we might be completing a provider name
		if !found && !strings.HasPrefix(currentToComplete, "/") {
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

		// 2. We are completing a path.
		// If not found (no colon), it's the default provider (onedrive).
		if !found {
			provider = fs.DefaultProviderPrefix
			path = currentToComplete
		}

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
