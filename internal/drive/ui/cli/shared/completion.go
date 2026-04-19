package shared

import (
	"context"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

type noOpAliasService struct{}

func newNoOpAliasService() *noOpAliasService {
	return &noOpAliasService{}
}

func (s *noOpAliasService) GetAliasByDriveID(ctx context.Context, driveID string) (string, error) {
	return "", nil
}

func ProviderPathCompletion(container di.Container, supportAlias bool) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var aliasSvc interface {
			GetAliasByDriveID(ctx context.Context, driveID string) (string, error)
		} = newNoOpAliasService()
		if supportAlias {
			aliasSvc = container.Alias()
		}

		drives, err := container.Drive().ListDrives(cmd.Context(), "")
		if err != nil {
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, drive := range drives {
			displayName := drive.Name
			alias, err := aliasSvc.GetAliasByDriveID(cmd.Context(), drive.Name)
			if err == nil && alias != "" {
				displayName = alias
			}

			if toComplete == "" || strings.HasPrefix(strings.ToLower(displayName), strings.ToLower(toComplete)) {
				results = append(results, displayName)
			}
		}
		return results, cobra.ShellCompDirectiveNoSpace | cobra.ShellCompDirectiveNoFileComp

	}
}
