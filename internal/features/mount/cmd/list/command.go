package list

import (
	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/spf13/cobra"
)

func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("mount-list")
	handler := NewCommand(container.Mounts(), formatting.NewFormatterFactory(), l)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List mount points",
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c = NewCommandContext(cmd.Context(), &opts)

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.Flags().StringVarP(&opts.Format, "format", "f", formatting.FormatShort.String(), "output format (short, long, json, yaml, tree, table)")

	return cmd
}
