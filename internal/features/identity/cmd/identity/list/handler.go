package list

import (
	"github.com/michaeldcanady/go-onedrive/pkg/format"
)

// IdentityListItem represents a single row in the identity list output.
type IdentityListItem struct {
	ID       string `json:"id" yaml:"id"`
	Email    string `json:"email" yaml:"email"`
	Provider string `json:"provider" yaml:"provider"`
}

// IdentityList is a collection of IdentityListItem that implements format.Tabular.
type IdentityList []IdentityListItem

// TableHeaders returns the headers for the table output.
func (l IdentityList) TableHeaders() []string {
	return []string{"ID", "EMAIL", "PROVIDER"}
}

// TableRows returns the rows for the table output.
func (l IdentityList) TableRows() [][]string {
	rows := make([][]string, len(l))
	for i, item := range l {
		rows[i] = []string{item.ID, item.Email, item.Provider}
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
	identities, err := c.identity.List(ctx.Ctx)
	if err != nil {
		return err
	}

	var list IdentityList
	for _, i := range identities {
		list = append(list, IdentityListItem{
			ID:       i.ID,
			Email:    i.Email,
			Provider: i.Provider,
		})
	}

	f := c.formatter.Get(format.Format(ctx.Options.Format))
	return f.Format(ctx.Options.Stdout, list)
}

// Finalize performs any cleanup or final output formatting.
func (c *Command) Finalize(ctx *CommandContext) error {
	return nil
}
