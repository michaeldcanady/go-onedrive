package list

import (
	"context"

	formatting "github.com/michaeldcanady/go-onedrive/pkg/format"
)

// CommandContext holds the execution state and configuration for the list operation.
type CommandContext struct {
	Ctx     context.Context
	Format  formatting.Format
	Options Options
}
