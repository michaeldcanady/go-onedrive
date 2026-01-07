package app

import (
	"context"
	"errors"
	"fmt"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
)

type GraphClientService struct {
	credentialService credentialLoader
	client            *msgraphsdkgo.GraphServiceClient
}

func NewGraphClientService(credSvc credentialLoader) *GraphClientService {
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

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(cred, []string{FilesReadWriteAllScope,
		UserReadScope})
	if err != nil {
		return nil, errors.Join(errors.New("unable to initialize client"), err)
	}

	s.client = client
	return client, nil
}
