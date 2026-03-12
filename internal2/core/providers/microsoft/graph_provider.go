package microsoft

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal2/core/logger"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

// GraphProvider facilitates the provisioning of an authenticated Microsoft Graph client.
type GraphProvider struct {
	// cred is the token credential used for authorizing requests.
	cred azcore.TokenCredential
	// log is the logger instance used for internal events.
	log logger.Logger
	// client is the cached Graph client instance.
	client *msgraphsdkgo.GraphServiceClient
}

// NewGraphProvider creates a new instance of GraphProvider with the provided credential and logger.
func NewGraphProvider(cred azcore.TokenCredential, log logger.Logger) *GraphProvider {
	return &GraphProvider{
		cred: cred,
		log:  log,
	}
}

// Client returns an authenticated Graph client, initializing it if necessary.
func (p *GraphProvider) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	if p.client != nil {
		return p.client, nil
	}

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(
		p.cred,
		[]string{
			"Files.ReadWrite.All",
			"User.Read",
			"offline_access",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create graph client: %w", err)
	}

	p.client = client
	return client, nil
}
