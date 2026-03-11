package domain

import "context"

type DriveAliasService interface {
	Resolve(ctx context.Context, alias string) (string, error)
}
