/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/app"
	"github.com/spf13/cobra"
)

var lsCmd = &cobra.Command{
	Use:   "ls [path]",
	Short: "List drives or items in a OneDrive path",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}

		path := ""
		if len(args) == 1 {
			path = args[0]
		}

		// Has CAE? - creates second cache
		driveSvc := app.NewDriveService(graphClientService)

		var cmdErr error
		driveSvc.ChildrenIterator(ctx, path)(func(name string, err error) bool {
			if err != nil {
				cmdErr = err
				return false
			}
			fmt.Println(name)
			return true
		})

		return cmdErr
	},
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
