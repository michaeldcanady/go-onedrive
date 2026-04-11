package cli

import (
	"path"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/spf13/cobra"
)

// ProviderPathCompletion provides shell completion for paths that may have a provider prefix (e.g., "local:/tmp").
// It handles the case where Bash splits words at colons by reconstructing the full path from args.
func ProviderPathCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {

		// Reconstruct full token in case Bash split on ':'
		full := strings.Join(append(args, toComplete), "")
		hasPrefix := strings.Contains(full, ":")

		uri, err := fs.ParseURI(full)
		if err != nil {
			// Probably completing provider name
			if !strings.Contains(full, ":") && !strings.HasPrefix(full, "/") {
				names, err := container.ProviderRegistry().RegisteredNames()
				if err != nil {
					return nil, cobra.ShellCompDirectiveError
				}
				var out []string
				for _, n := range names {
					if strings.HasPrefix(n, full) {
						out = append(out, n+":")
					}
				}
				return out, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
			}
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		// If path is empty → suggest root
		if uri.Path == "" || uri.Path == "." {
			root := "/"
			if hasPrefix {
				root = uri.Provider + ":" + root
			}
			return []string{root}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		// Split into directory + search prefix
		dir, file := path.Split(uri.Path)
		if dir == "" {
			dir = "/"
		}

		// Build manager path (provider:path)
		listPath := (&fs.URI{
			Provider: uri.Provider,
			DriveRef: uri.DriveRef,
			Path:     strings.TrimSuffix(dir, "/"),
		})

		items, err := container.FS().List(cmd.Context(), listPath, fs.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, item := range items {
			if !strings.HasPrefix(item.Name, file) {
				continue
			}

			entry := dir + item.Name
			if item.Type == fs.TypeFolder {
				entry += "/"
			}

			if hasPrefix {
				results = append(results, uri.Provider+":"+entry)
			} else {
				results = append(results, entry)
			}
		}

		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
