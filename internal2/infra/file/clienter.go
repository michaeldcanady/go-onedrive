package file

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type clienter interface {
	Client(context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
