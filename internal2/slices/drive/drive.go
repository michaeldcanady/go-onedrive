// Package drive provides the drive-related CLI commands.
package drive

import (
	"github.com/michaeldcanady/go-onedrive/internal2/di"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/alias"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/cat"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/cp"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/download"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/edit"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/get"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/list"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/ls"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/mkdir"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/mv"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/rm"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/touch"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/upload"
	"github.com/michaeldcanady/go-onedrive/internal2/slices/drive/use"
	"github.com/spf13/cobra"
)

// CreateDriveCmd constructs and returns the cobra.Command for the 'drive' parent command.
func CreateDriveCmd(container di.Container) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "drive",
		Short: "Manage OneDrive and local items",
	}

	cmd.AddCommand(
		ls.CreateLsCmd(container),
		cat.CreateCatCmd(container),
		mkdir.CreateMkdirCmd(container),
		rm.CreateRmCmd(container),
		touch.CreateTouchCmd(container),
		cp.CreateCpCmd(container),
		mv.CreateMvCmd(container),
		upload.CreateUploadCmd(container),
		download.CreateDownloadCmd(container),
		edit.CreateEditCmd(container),
		list.CreateListCmd(container),
		use.CreateUseCmd(container),
		get.CreateGetCmd(container),
		alias.CreateAliasCmd(container),
	)

	return cmd
}
