package ls

import (
	"encoding/json"
	"fmt"
	"os"

	"go.yaml.in/yaml/v3"
)

func printJSON(items []Item) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}

func printYAML(items []Item) error {
	out, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}

func printLong(items []Item) {
	for _, it := range items {
		kind := "-"
		if it.IsFolder {
			kind = "d"
		}

		fmt.Printf("%1s %10d %s %s\n",
			kind,
			it.Size,
			it.ModifiedTime.Format("2006-01-02 15:04"),
			it.Name,
		)
	}
}

func printShort(items []Item) {
	if len(items) == 0 {
		return
	}

	width := detectTerminalWidth()
	maxLen := 0
	for _, it := range items {
		if len(it.Name) > maxLen {
			maxLen = len(it.Name)
		}
	}

	colWidth := maxLen + 2
	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		fmt.Printf("%-*s", colWidth, it.Name)
		if (i+1)%cols == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}
