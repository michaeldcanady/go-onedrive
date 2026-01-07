package ls

import (
	"context"
	"iter"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

type driveChildIterator interface {
	ChildrenIterator(ctx context.Context, folderPath string) iter.Seq2[models.DriveItemable, error]
}
