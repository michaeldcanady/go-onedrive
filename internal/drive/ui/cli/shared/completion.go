package shared

import (
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/spf13/cobra"
)

type noOpStateService struct{}

func newNoOpStateService() *noOpStateService {
	return &noOpStateService{}
}

func (s *noOpStateService) GetAliasByDriveID(driveID string) (string, error) {
	return "", nil
}

func ProviderPathCompletion(container di.Container, supportAlias bool) func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var stateSvc interface {
			GetAliasByDriveID(driveID string) (string, error)
		} = newNoOpStateService()
		if supportAlias {
			stateSvc = container.Alias()
		}

		drives, err := container.Drive().ListDrives(cmd.Context())
		if err != nil {
			return nil, cobra.ShellCompDirectiveError | cobra.ShellCompDirectiveNoFileComp
		}

		var results []string
		for _, drive := range drives {
			displayName := drive.Name
			alias, err := stateSvc.GetAliasByDriveID(drive.Name)
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
