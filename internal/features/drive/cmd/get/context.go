package get

import (
	"context"
)

// CommandContext holds the execution state and configuration for the get operation.
type CommandContext struct {
	Ctx     context.Context
	Options Options
}
