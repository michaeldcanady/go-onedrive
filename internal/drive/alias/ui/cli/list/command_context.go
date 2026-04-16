package list

import (
	"context"
)

// CommandContext holds the execution state and configuration for the list operation.
type CommandContext struct {
	Ctx     context.Context
	Options Options
}
