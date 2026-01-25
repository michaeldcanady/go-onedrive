package cliprofileservicego

import "context"

type CacheService interface {
	GetCLIProfile(context.Context, string) (Profile, error)
	SetCLIProfile(context.Context, string, Profile) error
}
