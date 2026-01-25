package di

import (
	"context"

	cliprofileservicego "github.com/michaeldcanady/go-onedrive/internal/app/cli_profile_service.go"
)

type CLIProfileService interface {
	GetProfile(ctx context.Context, name string) (cliprofileservicego.Profile, error)
}
