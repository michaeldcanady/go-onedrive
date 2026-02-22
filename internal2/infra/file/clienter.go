package file

import (
	"context"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

// clienter is an internal interface defining the behavior for obtaining a Microsoft
// Graph Service Client. It facilitates dependency injection and mocking within
// the file package.
type clienter interface {
	// Client returns an initialized GraphServiceClient or an error if the client
	// cannot be created.
	Client(context.Context) (*msgraphsdkgo.GraphServiceClient, error)
}
