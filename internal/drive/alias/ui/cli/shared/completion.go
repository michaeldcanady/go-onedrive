package shared

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

func ProviderPathCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		aliases, err := container.Alias().ListAliases()
		if err != nil {
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, alias := range aliases {
			if toComplete == "" || strings.HasPrefix(strings.ToLower(alias), strings.ToLower(toComplete)) {
				results = append(results, alias)
			}
		}
		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp
	}
}
