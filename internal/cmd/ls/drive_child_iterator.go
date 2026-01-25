package ls

import (
	"context"
	"iter"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type driveChildIterator interface {
	Resolve(ctx context.Context, path string) (models.DriveItemable, error)
	ChildrenIterator(ctx context.Context, folderPath string) iter.Seq2[models.DriveItemable, error]
}
