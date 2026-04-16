package ls

import (
	"github.com/michaeldcanady/go-onedrive/internal/di"
	"github.com/michaeldcanady/go-onedrive/internal/fs/formatting"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/michaeldcanady/go-onedrive/pkg/args"
	"github.com/michaeldcanady/go-onedrive/pkg/flags"
	"github.com/spf13/cobra"
)

// CreateLsCmd constructs and returns the cobra.Command for the ls operation.
func CreateLsCmd(container di.Container) *cobra.Command {
	var opts Options
	var c *CommandContext

	l, _ := container.Logger().CreateLogger("ls")
	handler := NewCommand(container.FS(), container.URIFactory(), formatting.NewFormatterFactory(), l)

	cmd := &cobra.Command{
		Use:               "ls <path>",
		Short:             "List items in a directory",
		Long:              "List the items in a specified directory in OneDrive or the local filesystem.",
		Args:              args.MaximumNArgs(&opts),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, argsSlice []string) error {
			if err := args.Bind(argsSlice, &opts); err != nil {
				return err
			}
			opts.Stdout = cmd.OutOrStdout()
			opts.Format = formatting.NewFormat(opts.FormatStr)

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
		panic(err) // Registration should not fail if struct is correct
	}

	return cmd
}
