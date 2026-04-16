package cp

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/pkg/args"
	"github.com/michaeldcanady/go-onedrive/pkg/flags"
	"github.com/spf13/cobra"
)

// CreateCpCmd constructs and returns the cobra.Command for the drive cp operation.
func CreateCpCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("drive-cp")
	handler := NewCommand(container.FS(), container.URIFactory(), l)

	cmd := &cobra.Command{
		Use:               "cp <source> <destination>",
		Short:             "Copy files and directories",
		Args:              args.ExactArgs(&opts),
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
