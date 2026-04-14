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

		// 1. Determine if we are completing a provider/alias or a path
		uri, err := container.URIFactory().FromString(currentToComplete)
		found := strings.Contains(currentToComplete, ":")

		// If parsing failed and it doesn't look like a path (no slash), we might be completing a provider/alias name
		if err != nil && !strings.HasPrefix(currentToComplete, "/") && !strings.Contains(currentToComplete, "/") {
			names, _ := container.ProviderRegistry().RegisteredNames()
			// TODO: Add aliases to completion?

			var results []string
			for _, name := range names {
				if strings.HasPrefix(name, currentToComplete) {
					results = append(results, name+":")
				}
			}
			return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		if err != nil {
			// If it has a colon but isn't a known provider/alias, we can't complete it easily here.
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// 2. We are completing a path within a provider/alias.
		// If path is empty (just "provider:"), suggest "/"
		if uri.Path == "" {
			res := "/"
			if found && prefixToStrip == "" {
				res = currentToComplete + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		// Find parent directory and prefix
		lastSlash := strings.LastIndex(uri.Path, "/")
		if lastSlash == -1 {
			// No slash after the colon (e.g. "local:tm") or path is just "tm"
			res := "/"
			if found && prefixToStrip == "" {
				// Reconstruct the "provider:" part
				prefix, _, _ := strings.Cut(currentToComplete, ":")
				res = prefix + ":" + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		listDirPath := uri.Path[:lastSlash+1]
		searchPrefix := uri.Path[lastSlash+1:]

		// Trim trailing slash for validation (except for root "/")
		cleanListDirPath := listDirPath
		if len(cleanListDirPath) > 1 && strings.HasSuffix(cleanListDirPath, "/") {
			cleanListDirPath = strings.TrimSuffix(cleanListDirPath, "/")
		}

		// Create a URI for the directory to list
		listURI := *uri
		listURI.Path = cleanListDirPath

		items, err := container.FS().List(cmd.Context(), &listURI, fs.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, item := range items {
			if strings.HasPrefix(item.Name, searchPrefix) {
				res := listDirPath + item.Name
				if item.Type == fs.TypeFolder {
					res += "/"
				}

				// If we reconstructed the prefix from 'args' (meaning the shell split it),
				// we return the path without the prefix so the shell matches it correctly.
				// If the prefix was in 'toComplete', we return the full path.
				if found && prefixToStrip == "" {
					prefix, _, _ := strings.Cut(currentToComplete, ":")
					results = append(results, prefix+":"+res)
				} else {
					results = append(results, res)
				}
			}
		}

		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
