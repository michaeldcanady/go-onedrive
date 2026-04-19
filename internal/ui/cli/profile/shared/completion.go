package shared

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

func ProviderPathCompletion(container di.Container) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		profiles, err := container.Profile().List(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, profile := range profiles {

			if toComplete == "" || strings.HasPrefix(strings.ToLower(profile.Name), strings.ToLower(toComplete)) {
				results = append(results, profile.Name)
			}
		}
		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp

	}
}
