package remove

import (
	"context"
)

// CommandContext holds the execution state and configuration for the remove operation.
type CommandContext struct {
	Ctx     context.Context
	Options Options
}
