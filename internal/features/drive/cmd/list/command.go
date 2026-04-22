package list

import (
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/di"
	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
	"github.com/spf13/cobra"
)

// CreateListCmd constructs and returns the cobra.Command for the drive list operation.
func CreateListCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-list")
	handler := NewCommand(container.Drive(), formatting.NewFormatterFactory(), l)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available OneDrive drives",
		Long:  `Retrieve all OneDrive drives associated with your account, marking the personal drive with an asterisk.`,
		Args:  cobra.NoArgs,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Stdout = cmd.OutOrStdout()
			opts.Stderr = cmd.ErrOrStderr()

			c = &CommandContext{
				Ctx:     cmd.Context(),
				Options: opts,
			}

			return handler.Validate(c)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := handler.Execute(c); err != nil {
				return err
			}
			return handler.Finalize(c)
		},
	}

	cmd.Flags().StringVar(&opts.IdentityID, "id", "", "The specific identity (email or alias) to list drives for")
	cmd.Flags().StringVarP(&opts.Format, "format", "o", "table", fmt.Sprintf("Output format %s", supportedFormats))

	return cmd
}
