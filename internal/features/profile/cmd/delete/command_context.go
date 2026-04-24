package delete

import (
	"context"
)

type CommandContext struct {
	Ctx     context.Context
	Options Options
}
