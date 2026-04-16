package current

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal/profile"
)

type CommandContext struct {
	Ctx     context.Context
	Options Options

	Profile *profile.Profile
}
