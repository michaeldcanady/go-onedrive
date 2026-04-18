// Package cli provides the config-related CLI commands.
package config

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/alias/ui/cli/set"
	"github.com/michaeldcanady/go-onedrive/internal/drive/ui/cli/get"
	"github.com/spf13/cobra"
)

// CreateConfigCmd constructs and returns the cobra.Command for the 'config' parent command.
func CreateConfigCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration settings",
	}

	cmd.AddCommand(
		get.CreateGetCmd(container),
		set.CreateSetCmd(container),
	)

	return cmd
}
