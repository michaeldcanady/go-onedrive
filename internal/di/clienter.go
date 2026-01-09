package di

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type Clienter interface {
	Client(context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
