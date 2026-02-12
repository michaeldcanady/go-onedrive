package profile

import "context"

type ProfileService interface {
	Get(ctx context.Context, name string) (Profile, error)
	List(ctx context.Context) ([]Profile, error)
	Create(ctx context.Context, name string) (Profile, error)
	Delete(ctx context.Context, name string) error
	Exists(ctx context.Context, name string) (bool, error)
}
