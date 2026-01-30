package graph

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"

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
	logger            logging.Logger
}

func New(credSvc azcore.TokenCredential, logger logging.Logger) *GraphClientService {
	return &GraphClientService{
		credentialService: credSvc,
		logger:            logger,
	}
}

func (s *GraphClientService) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	// Already initialized?
	if s.client != nil {
		s.logger.Debug("graph client already initialized")
		return s.client, nil
	}

	s.logger.Info("credential loaded successfully for graph client")

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(
		s.credentialService,
		[]string{
			FilesReadWriteAllScope,
			UserReadScope,
			//OfflineAccessScope,
		},
	)
	if err != nil {
		s.logger.Error("unable to initialize graph client", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to initialize client"), err)
	}

	s.client = client

	s.logger.Info("graph client initialized successfully")
	s.logger.Debug("graph client instance", logging.Any("client", client))

	return client, nil
}
