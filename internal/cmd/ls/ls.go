package ls

import (
	"cmp"
	"context"
	"os"
	"slices"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

func CreateLSCmd(driveChildIterator driveChildIterator) *cobra.Command {
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

			var cmdErr error
			items := make([]string, 0)

			// Collect items
			driveChildIterator.ChildrenIterator(ctx, path)(func(item models.DriveItemable, err error) bool {
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

				items = appendSorted(items, name)
				return true
			})

			if cmdErr != nil {
				return cmdErr
			}

			printColumns2(items)

			return nil

		},
	}

	return lsCmd
}

// Source - https://stackoverflow.com/a
// Posted by Andrew W. Phillips, modified by community. See post 'Timeline' for change history
// Retrieved 2026-01-01, License - CC BY-SA 4.0

func appendSorted[T cmp.Ordered](ts []T, t T) []T {
	i, _ := slices.BinarySearch(ts, t)
	return slices.Insert(ts, i, t)
}

func printColumns2(names []string) {
	if len(names) == 0 {
		return
	}

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

	colWidth := maxLen + 2
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	// Chunk into rows
	rows := make([][]string, 0)
	for i := 0; i < len(names); i += cols {
		end := i + cols
		if end > len(names) {
			end = len(names)
		}

		row := names[i:end]

		// pad row to full width
		for len(row) < cols {
			row = append(row, "")
		}

		rows = append(rows, row)
	}

	symbols := tw.NewSymbolCustom("linux").
		WithRow(" ").
		WithColumn(" ").
		WithTopLeft(" ").
		WithTopMid(" ").
		WithTopRight(" ").
		WithMidLeft(" ").
		WithCenter(" ").
		WithMidRight(" ").
		WithBottomLeft(" ").
		WithBottomMid(" ").
		WithBottomRight(" ")

	// Render table
	table := tablewriter.NewTable(os.Stdout, tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: symbols})))
	table.Header([]string{})

	for _, row := range rows {
		table.Append(row)
	}

	table.Render()
}
