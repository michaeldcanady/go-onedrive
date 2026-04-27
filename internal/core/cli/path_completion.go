package cli

import (
	"context"
	"strings"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
	"github.com/michaeldcanady/go-onedrive/internal/features/mount"
	pkgfs "github.com/michaeldcanady/go-onedrive/pkg/fs"
	"github.com/spf13/cobra"
)

type uriFactory interface {
	FromString(input string) (*fs.URI, error)
}

type mountLister interface {
	ListMounts(ctx context.Context) ([]mount.MountConfig, error)
}

type itemLister interface {
	List(ctx context.Context, uri *pkgfs.URI, opts pkgfs.ListOptions) ([]pkgfs.Item, error)
}

// ProviderPathCompletion provides shell completion for paths that may have a provider prefix (e.g., "local:/tmp").
func ProviderPathCompletion(itemLister itemLister, uriFactory uriFactory, mountLister mountLister) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		currentToComplete, prefixToStrip := resolveToComplete(args, toComplete)

		uri, err := uriFactory.FromString(currentToComplete)
		found := strings.Contains(currentToComplete, ":")

		if err != nil && !strings.ContainsAny(currentToComplete, "/:") {
			return mountCompletion(cmd.Context(), mountLister, currentToComplete)
		}

		if err != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

		if uri.Path == "" {
			return rootPathCompletion(currentToComplete, found, prefixToStrip), cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		lastSlash := strings.LastIndex(uri.Path, "/")
		if lastSlash == -1 {
			return []string{formatPathResult(currentToComplete, found, prefixToStrip, "/")}, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
		}

		return listPathCompletion(cmd.Context(), itemLister, uri, lastSlash)
	}
}

func resolveToComplete(args []string, toComplete string) (string, string) {
	currentToComplete := toComplete
	prefixToStrip := ""
	if !strings.Contains(toComplete, ":") && len(args) > 0 {
		lastArg := args[len(args)-1]
		if lastArg == ":" && len(args) >= 2 {
			prefixToStrip = args[len(args)-2] + ":"
			currentToComplete = prefixToStrip + toComplete
		} else if strings.HasSuffix(lastArg, ":") {
			prefixToStrip = lastArg
			currentToComplete = prefixToStrip + toComplete
		}
	}
	return currentToComplete, prefixToStrip
}

func mountCompletion(ctx context.Context, mountLister mountLister, toComplete string) ([]string, cobra.ShellCompDirective) {
	mounts, _ := mountLister.ListMounts(ctx)
	var results []string
	for _, m := range mounts {
		name := strings.TrimPrefix(m.Path, "/")
		if name != "" && strings.HasPrefix(name, toComplete) {
			results = append(results, name+":")
		}
	}
	return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
}

func rootPathCompletion(currentToComplete string, found bool, prefixToStrip string) []string {
	res := "/"
	if found && prefixToStrip == "" {
		res = currentToComplete + res
	}
	return []string{res}
}

func formatPathResult(currentToComplete string, found bool, prefixToStrip string, path string) string {
	if found && prefixToStrip == "" {
		prefix, _, _ := strings.Cut(currentToComplete, ":")
		return prefix + ":" + path
	}
	return path
}

func listPathCompletion(ctx context.Context, itemLister itemLister, uri *pkgfs.URI, lastSlash int) ([]string, cobra.ShellCompDirective) {
	listDirPath := uri.Path[:lastSlash+1]
	searchPrefix := uri.Path[lastSlash+1:]

	cleanListDirPath := listDirPath
	if len(cleanListDirPath) > 1 && strings.HasSuffix(cleanListDirPath, "/") {
		cleanListDirPath = strings.TrimSuffix(cleanListDirPath, "/")
	}

	listURI := *uri
	listURI.Path = cleanListDirPath

	items, err := itemLister.List(ctx, &listURI, pkgfs.ListOptions{})
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
