package list

import (
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// MountListItem represents a single row in the mount list output.
type MountListItem struct {
	Path     string `json:"path" yaml:"path"`
	Type     string `json:"type" yaml:"type"`
	Identity string `json:"identity" yaml:"identity"`
}

// MountList is a collection of MountListItem that implements format.Tabular.
type MountList []MountListItem

// TableHeaders returns the headers for the table output.
func (l MountList) TableHeaders() []string {
	return []string{"PATH", "TYPE", "IDENTITY"}
}

// TableRows returns the rows for the table output.
func (l MountList) TableRows() [][]string {
	rows := make([][]string, len(l))
	for i, item := range l {
		rows[i] = []string{item.Path, item.Type, item.Identity}
	}
	return rows
}

// Validate performs initial validation of the command options.
func (c *Command) Validate(ctx *CommandContext) error {
	return nil
}

// Resolve performs argument resolution.
func (c *Command) Resolve(ctx *CommandContext) error {
	return c.BaseResolve(ctx)
}

// Execute performs the core business logic of the command.
func (c *Command) Execute(ctx *CommandContext) error {
	mounts, err := c.mounts.List(ctx.Ctx)
	if err != nil {
		return err
	}

	var list MountList
	for _, m := range mounts {
		list = append(list, MountListItem{
			Path:     m.Path,
			Type:     m.Type,
			Identity: m.IdentityID,
		})
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, list)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
