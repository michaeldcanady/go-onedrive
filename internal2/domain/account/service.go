package account

import (
	"context"
)

type Service interface {
	Get(context.Context) (Account, error)
	Put(context.Context, Account) error
	Delete(context.Context) error
}
