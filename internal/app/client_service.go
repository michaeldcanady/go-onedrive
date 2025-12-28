package app

import (
	"context"
	"errors"
	"fmt"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

type GraphClientService struct {
	credentialService CredentialService
	client            *msgraphsdkgo.GraphServiceClient
}

func NewGraphClientService(credSvc CredentialService) *GraphClientService {
	return &GraphClientService{
		credentialService: credSvc,
	}
}

func (s *GraphClientService) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	// Already initialized?
	if s.client != nil {
		return s.client, nil
	}

	cred, err := s.credentialService.LoadCredential(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to load credential: %w", err)
	}

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{"Files.ReadWrite"})
	if err != nil {
		return nil, errors.Join(errors.New("unable to initialize client"), err)
	}

	s.client = client
	return client, nil
}
