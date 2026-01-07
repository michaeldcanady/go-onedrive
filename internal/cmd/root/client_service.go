package root

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type ClientService interface {
	Client(context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
