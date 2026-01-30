package profile

import "context"

type ProfileService interface {
	GetProfile(ctx context.Context, name string) (Profile, error)
}
