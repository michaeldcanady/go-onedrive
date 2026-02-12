package profile

import "context"

type Repository interface {
	Get(ctx context.Context, name string) (Profile, error)
	List() ([]Profile, error)
	Create(name string) (Profile, error)
	Delete(name string) error
	Exists(name string) (bool, error)
}
