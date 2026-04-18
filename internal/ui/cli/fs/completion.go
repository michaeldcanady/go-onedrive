package fs

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/spf13/cobra"
)

// ProviderPathCompletion provides shell completion for paths that may have a provider prefix (e.g., "local:/tmp").
func ProviderPathCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		currentToComplete := toComplete
		prefixToStrip := ""

		if !strings.Contains(toComplete, ":") && len(args) > 0 {
			if args[len(args)-1] == ":" && len(args) >= 2 {
				prefixToStrip = args[len(args)-2] + ":"
				currentToComplete = prefixToStrip + toComplete
			} else if strings.HasSuffix(args[len(args)-1], ":") {
				prefixToStrip = args[len(args)-1]
				currentToComplete = prefixToStrip + toComplete
			}
		}

		uri, err := container.URIFactory().FromString(currentToComplete)
		found := strings.Contains(currentToComplete, ":")

		if err != nil && !strings.HasPrefix(currentToComplete, "/") && !strings.Contains(currentToComplete, "/") {
			names, _ := container.ProviderRegistry().RegisteredNames()
			var results []string
			for _, name := range names {
				if strings.HasPrefix(name, currentToComplete) {
					results = append(results, name+":")
				}
			}
			return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if uri.Path == "" {
			res := "/"
			if found && prefixToStrip == "" {
				res = currentToComplete + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		lastSlash := strings.LastIndex(uri.Path, "/")
		if lastSlash == -1 {
			res := "/"
			if found && prefixToStrip == "" {
				prefix, _, _ := strings.Cut(currentToComplete, ":")
				res = prefix + ":" + res
			}
			return []string{res}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		listDirPath := uri.Path[:lastSlash+1]
		searchPrefix := uri.Path[lastSlash+1:]

		cleanListDirPath := listDirPath
		if len(cleanListDirPath) > 1 && strings.HasSuffix(cleanListDirPath, "/") {
			cleanListDirPath = strings.TrimSuffix(cleanListDirPath, "/")
		}

		listURI := *uri
		listURI.Path = cleanListDirPath

		items, err := container.FS().List(cmd.Context(), &listURI, pkgfs.ListOptions{})
		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, item := range items {
			if strings.HasPrefix(item.Name, searchPrefix) {
				res := listDirPath + item.Name
				if item.Type == pkgfs.TypeFolder {
					res += "/"
				}
				results = append(results, res)
			}
		}
		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
