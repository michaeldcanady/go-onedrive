package list

import (
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// DriveListItem represents a single row in the drive list output.
type DriveListItem struct {
	Mounted  string `json:"mounted" yaml:"mounted"`
	ID       string `json:"id" yaml:"id"`
	Name     string `json:"name" yaml:"name"`
	Type     string `json:"type" yaml:"type"`
	Identity string `json:"identity" yaml:"identity"`
}

// DriveList is a collection of DriveListItem that implements format.Tabular.
type DriveList []DriveListItem

// TableHeaders returns the headers for the table output.
func (l DriveList) TableHeaders() []string {
	return []string{"MOUNTED", "ID", "NAME", "TYPE", "IDENTITY"}
}

// TableRows returns the rows for the table output.
func (l DriveList) TableRows() [][]string {
	rows := make([][]string, len(l))
	for i, item := range l {
		rows[i] = []string{item.Mounted, item.ID, item.Name, item.Type, item.Identity}
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
	identityID := ctx.Options.Id
	drives, err := c.drive.List(ctx.Ctx, identityID)
	if err != nil {
		return err
	}

	mounts, err := c.mount.List(ctx.Ctx)
	if err != nil {
		return err
	}

	var list DriveList
	for _, d := range drives {
		mountedPath := ""
		for _, m := range mounts {
			if m.IdentityID == d.IdentityID { // Simplified check
				mountedPath = m.Path
				break
			}
		}
		list = append(list, DriveListItem{
			Mounted:  mountedPath,
			ID:       d.ID,
			Name:     d.Name,
			Type:     d.Type,
			Identity: d.IdentityID,
		})
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, list)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
