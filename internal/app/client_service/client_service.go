package clientservice

import (
	"context"
	"errors"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	msgraphsdkgo "github.com/microsoftgraph/msgraph-sdk-go"
)

const (
	FilesReadWriteAllScope = "Files.ReadWrite.All"
	UserReadScope          = "User.Read"
)

type GraphClientService struct {
	credentialService credentialLoader
	client            *msgraphsdkgo.GraphServiceClient
	publisher         event.Publisher
	logger            logging.Logger
}

func New(credSvc credentialLoader, publisher event.Publisher, logger logging.Logger) *GraphClientService {
	return &GraphClientService{
		credentialService: credSvc,
		publisher:         publisher,
		logger:            logger,
	}
}

func (s *GraphClientService) Client(ctx context.Context) (*msgraphsdkgo.GraphServiceClient, error) {
	// Already initialized?
	if s.client != nil {
		s.logger.Debug("graph client already initialized")
		return s.client, nil
	}

	cred, err := s.credentialService.LoadCredential(ctx)
	if err != nil {
		s.logger.Error("failed to load credential", logging.Any("error", err))
		return nil, fmt.Errorf("failed to load credential: %w", err)
	}

	s.logger.Info("credential loaded successfully for graph client")

	client, err := msgraphsdkgo.NewGraphServiceClientWithCredentials(
		cred,
		[]string{
			FilesReadWriteAllScope,
			UserReadScope,
		},
	)
	if err != nil {
		s.logger.Error("unable to initialize graph client", logging.Any("error", err))
		return nil, errors.Join(errors.New("unable to initialize client"), err)
	}

	s.client = client

	s.logger.Info("graph client initialized successfully")
	s.logger.Debug("graph client instance", logging.Any("client", client))

	// 🔥 Publish event
	_ = s.publisher.Publish(ctx, newGraphClientInitializedEvent())

	return client, nil
}
