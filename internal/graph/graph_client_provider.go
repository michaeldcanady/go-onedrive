package graph

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/auth"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphClientProvider struct {
	creds  auth.CredentialProvider
	client *msgraphsdkgo.GraphServiceClient
	log    logging.Logger
}

func NewGraphClientProvider(
	creds auth.CredentialProvider,
	log logging.Logger,
) *GraphClientProvider {
	return &GraphClientProvider{
		creds: creds,
		log:   log,
	}
}

func (p *GraphClientProvider) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	if p.client != nil {
		return p.client, nil
	}

	cred, err := p.creds.Credential(ctx, "default")
	if err != nil {
		return nil, fmt.Errorf("load credential: %w", err)
	}

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(
		cred,
		[]string{
			"Files.ReadWrite.All",
			"User.Read",
		},
	)
	if err != nil {
		return nil, fmt.Errorf("init graph client: %w", err)
	}

	p.client = client
	return client, nil
}
