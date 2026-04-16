package use

import (
	"context"
)

// CommandContext holds the execution state and configuration for the use operation.
type CommandContext struct {
	Ctx     context.Context
	Options Options
}
