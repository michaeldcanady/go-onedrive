package domain

import (
	"context"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type ClientProvider interface {
	Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error)
	RequestAdapter(ctx context.Context) (abstractions.RequestAdapter, error)
}
