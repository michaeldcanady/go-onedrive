package graph

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
	OfflineAccessScope     = "offline_access"
)

type GraphClientService struct {
	credentialService azcore.TokenCredential
	client            *msgraphsdkgo.GraphServiceClient
	log               logger.Logger
}

func New(credSvc azcore.TokenCredential, l logger.Logger) *GraphClientService {
	return &GraphClientService{
		credentialService: credSvc,
		log:               l,
	}
}

func (s *GraphClientService) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	// Already initialized?
	if s.client != nil {
		s.log.Debug("graph client already initialized")
		return s.client, nil
	}

	s.log.Info("credential loaded successfully for graph client")

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(
		s.credentialService,
		[]string{
			FilesReadWriteAllScope,
			UserReadScope,
			//OfflineAccessScope,
		},
	)
	if err != nil {
		s.log.Error("unable to initialize graph client", logger.Any("error", err))
		return nil, errors.Join(errors.New("unable to initialize client"), err)
	}

	s.client = client

	s.log.Info("graph client initialized successfully")
	s.log.Debug("graph client instance", logger.Any("client", client))

	return client, nil
}

func (s *GraphClientService) RequestAdapter(ctx context.Context) (abstractions.RequestAdapter, error) {
	client, err := s.Client(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetAdapter(), nil
}
