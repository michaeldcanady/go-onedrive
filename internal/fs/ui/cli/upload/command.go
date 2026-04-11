package upload

import (
	"errors"
	"fmt"
	"slices"

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
			// 1. Process Source Path
			srcURI, err := fs.ParseURI(args[0])
			if err != nil {
				return coreerrors.NewInvalidInput(
					err,
					fmt.Sprintf("invalid source path '%s'", args[0]),
					"Check the path format and ensure no illegal characters are used.",
				)
			}
			opts.Source = srcURI

			if err := fs.ValidatePathSyntax(srcURI.Path); err != nil {
				if _, ok := err.(*coreerrors.TrailingSlashError); ok && !opts.Recursive {
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid source path '%s' due to trailing slash", opts.Source.String()),
						"Remove the trailing slash from the source path",
					)
				}
				if !opts.Recursive && err != nil {
					return err
				}
			}

			// 2. Process Destination Path
			dstURI, err := fs.ParseURI(args[1])
			if err != nil {
				return coreerrors.NewInvalidInput(
					err,
					fmt.Sprintf("invalid destination path '%s'", args[1]),
					"Check the path format and ensure no illegal characters are used.",
				)
			}
			opts.Destination = dstURI

			if err := fs.ValidatePathSyntax(dstURI.Path); err != nil {
				if _, ok := err.(*coreerrors.TrailingSlashError); ok && !opts.Recursive {
					return coreerrors.NewInvalidInput(
						err,
						fmt.Sprintf("invalid destination path '%s' due to trailing slash", opts.Destination.String()),
						"Remove the trailing slash from the destination path",
					)
				}
				if !opts.Recursive && err != nil {
					return err
				}
			}

			// 3. Check Providers
			names, err := container.ProviderRegistry().RegisteredNames()
			if err != nil {
				return coreerrors.NewAppError(
					coreerrors.CodeUnknown,
					errors.New("failed to check registered providers"),
					"An unexpected error occurred while retrieving registered providers",
					"Try again, and if the problem persists, check the application logs for more details",
				)
			}

			if !slices.Contains(names, srcURI.Provider) {
				return coreerrors.NewInvalidInput(
					errors.New("unknown source provider"),
					fmt.Sprintf("unknown source provider prefix '%s'", srcURI.Provider),
					"Ensure the provider prefix is correct and corresponds to a registered provider",
				)
			}

			if !slices.Contains(names, dstURI.Provider) {
				return coreerrors.NewInvalidInput(
					errors.New("unknown destination provider"),
					fmt.Sprintf("unknown destination provider prefix '%s'", dstURI.Provider),
					"Ensure the provider prefix is correct and corresponds to a registered provider",
				)
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
