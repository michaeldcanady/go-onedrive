package cp

import "context"

type CommandContext struct {
	Ctx     context.Context
	Options Options
}
