package list

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/features/profile"
)

type CommandContext struct {
	Ctx     context.Context
	Options Options

	// Execution results
	Profiles []profile.Profile
	Active   *profile.Profile
}
