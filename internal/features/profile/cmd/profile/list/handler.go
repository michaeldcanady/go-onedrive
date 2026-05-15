package list

import (
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// ProfileListItem represents a single row in the profile list output.
type ProfileListItem struct {
	Active bool   `json:"active" yaml:"active"`
	Name   string `json:"name" yaml:"name"`
}

// ProfileList is a collection of ProfileListItem that implements format.Tabular.
type ProfileList []ProfileListItem

// TableHeaders returns the headers for the table output.
func (l ProfileList) TableHeaders() []string {
	return []string{"ACTIVE", "NAME"}
}

// TableRows returns the rows for the table output.
func (l ProfileList) TableRows() [][]string {
	rows := make([][]string, len(l))
	for i, item := range l {
		active := ""
		if item.Active {
			active = "*"
		}
		rows[i] = []string{active, item.Name}
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
	profiles, err := c.profile.List()
	if err != nil {
		return err
	}

	current, err := c.profile.GetCurrent()
	if err != nil {
		// Ignore error if no current profile
		current = nil
	}

	var list ProfileList
	for _, p := range profiles {
		isActive := current != nil && current.Name == p.Name
		list = append(list, ProfileListItem{
			Active: isActive,
			Name:   p.Name,
		})
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, list)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
