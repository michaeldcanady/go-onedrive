package di

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphClientProvider interface {
	Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
