package graph

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type ClientProvider interface {
	Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
