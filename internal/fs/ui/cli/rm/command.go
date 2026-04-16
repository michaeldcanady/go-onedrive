package rm

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/pkg/args"
	"github.com/michaeldcanady/go-onedrive/pkg/flags"
	"github.com/spf13/cobra"
)

// CreateRmCmd constructs and returns the cobra.Command for the drive rm operation.
func CreateRmCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-rm")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "rm <path>",
		Short:             "Remove a file or directory",
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, argsSlice []string) error {
			if err := args.Bind(argsSlice, &opts); err != nil {
				return err
			}
			opts.Stdout = cmd.OutOrStdout()

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

	if err := flags.RegisterFlags(cmd, &opts); err != nil {
		panic(err)
	}

	return cmd
}
