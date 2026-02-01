package cliprofileservicego

import (
	"context"

	domainprofile "github.com/michaeldcanady/go-onedrive/internal2/domain/profile"
)

type CacheService interface {
	GetCLIProfile(context.Context, string) (domainprofile.Profile, error)
	SetCLIProfile(context.Context, string, domainprofile.Profile) error
}
