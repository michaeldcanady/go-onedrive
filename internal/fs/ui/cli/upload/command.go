package upload

import (
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/di"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
	"github.com/michaeldcanady/go-onedrive/internal/fs"
	"github.com/michaeldcanady/go-onedrive/internal/fs/ui/cli"
	"github.com/spf13/cobra"
)

// CreateUploadCmd constructs and returns the cobra.Command for the drive upload operation.
func CreateUploadCmd(container di.Container) *cobra.Command {
	var opts Options

	cmd := &cobra.Command{
		Use:               "upload <local_path> <remote_path>",
		Short:             "Upload files and directories to OneDrive",
		Args:              cobra.ExactArgs(2),
		ValidArgsFunction: cli.ProviderPathCompletion(container),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			opts.Source = args[0]
			if opts.Recursive && fs.ContainsIllegalChars(opts.Source) {
				return coreerrors.NewInvalidInput(
					errors.New("source path contains illegal characters"),
					fmt.Sprintf("invalid source path '%s' due to illegal characters", opts.Source),
					"Remove the illegal characters from the source path",
				)
			} else {
				if err := fs.ValidatePathSyntax(opts.Source); err != nil {
					return err
				}
			}

			opts.Destination = args[1]
			if opts.Recursive && fs.ContainsIllegalChars(opts.Destination) {
				return coreerrors.NewInvalidInput(
					errors.New("destination path contains illegal characters"),
					fmt.Sprintf("invalid destination path '%s' due to illegal characters", opts.Destination),
					"Remove the illegal characters from the destination path",
				)
			} else {
				if err := fs.ValidatePathSyntax(opts.Destination); err != nil {
					return err
				}
			}

			opts.Stdout = cmd.OutOrStdout()
			return opts.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return NewHandler(container.FS(), container.Logger()).Handle(cmd.Context(), opts)
		},
	}

	cmd.Flags().BoolVarP(&opts.Recursive, "recursive", "r", false, "upload directories recursively")

	return cmd
}
