package di

import (
	"context"
	"iter"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type ChildrenIterator interface {
	ChildrenIterator(context.Context, string) iter.Seq2[models.DriveItemable, error]
	Resolve(ctx context.Context, path string) (models.DriveItemable, error)
}
