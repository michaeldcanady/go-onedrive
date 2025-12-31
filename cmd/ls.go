/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/app"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/spf13/cobra"
	"golang.org/x/term"
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

		driveSvc := app.NewDriveService(graphClientService)

		var cmdErr error
		items := make([]string, 0)

		// Collect items
		driveSvc.ChildrenIterator(ctx, path)(func(item models.DriveItemable, err error) bool {
			if err != nil {
				cmdErr = err
				return false
			}
			name := ""
			if item.GetName() != nil {
				name = *item.GetName()
			}

			if item.GetFolder() != nil {
				name = name + "/"
			}

			items = append(items, name)
			return true
		})

		if cmdErr != nil {
			return cmdErr
		}

		// Sort alphabetically like Linux ls
		sort.Slice(items, func(i, j int) bool {
			nameI := items[i]
			nameJ := items[j]

			return strings.ToLower(nameI) < strings.ToLower(nameJ)
		})

		printColumns(items)

		return nil

	},
}

func printColumns(names []string) {
	// Determine terminal width
	width := 80 // fallback
	if w, _, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
	}

	// Determine max item width
	maxLen := 0
	for _, n := range names {
		if len(n) > maxLen {
			maxLen = len(n)
		}
	}

	// Add padding between columns
	colWidth := maxLen + 2

	// Determine how many columns fit
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	// Print rows
	for i := 0; i < len(names); i += cols {
		end := i + cols
		if end > len(names) {
			end = len(names)
		}

		row := names[i:end]
		for _, n := range row {
			fmt.Printf("%-*s", colWidth, n)
		}
		fmt.Println()
	}
}

func init() {
	rootCmd.AddCommand(lsCmd)
}
