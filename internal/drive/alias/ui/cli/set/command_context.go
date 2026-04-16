package set

import (
	"context"
)

// CommandContext holds the execution state and configuration for the set operation.
type CommandContext struct {
	Ctx     context.Context
	Options Options
}
