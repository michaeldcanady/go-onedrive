package ls

import (
	"encoding/json"
	"fmt"
	"os"

	driveservice "github.com/michaeldcanady/go-onedrive/internal/app/drive_service"
	"go.yaml.in/yaml/v3"
)

func printJSON(items []*driveservice.DriveItem) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(items)
}

func printYAML(items []*driveservice.DriveItem) error {
	out, err := yaml.Marshal(items)
	if err != nil {
		return err
	}
	fmt.Print(string(out))
	return nil
}

func printLongDomain(items []*driveservice.DriveItem) {
	for _, it := range items {
		mod := it.Modified.Format("2006-01-02 15:04")

		size := "-"
		if !it.IsFolder {
			size = fmt.Sprintf("%d", it.Size)
		}

		name := it.Name
		if it.IsFolder {
			name += "/"
		}

		fmt.Printf("%-20s %10s  %s\n", mod, size, name)
	}
}

func printShortDomain(items []*driveservice.DriveItem) {
	width := detectTerminalWidth()
	colWidth := 0

	// Determine widest name
	for _, it := range items {
		if len(it.Name) > colWidth {
			colWidth = len(it.Name)
		}
	}
	colWidth += 2 // spacing

	cols := width / colWidth
	if cols < 1 {
		cols = 1
	}

	for i, it := range items {
		name := it.Name
		if it.IsFolder {
			name += "/"
		}

		fmt.Printf("%-*s", colWidth, name)

		if (i+1)%cols == 0 {
			fmt.Println()
		}
	}

	fmt.Println()
}
