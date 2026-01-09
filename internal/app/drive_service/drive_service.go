package driveservice

import (
	"context"
	"errors"
	"fmt"
	"iter"
	"path"

	"github.com/michaeldcanady/go-onedrive/internal/event"
	"github.com/michaeldcanady/go-onedrive/internal/logging"
	msgraphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
	"github.com/microsoftgraph/msgraph-sdk-go/drives"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

const (
	rootChildrenURITemplate         = "https://graph.microsoft.com/v1.0/drives/%s/root/children"
	rootRelativeChildrenURITemplate = "https://graph.microsoft.com/v1.0/drives/%s/root:%s:/children"
)

type Service struct {
	graph     clienter
	publisher event.Publisher
	logger    logging.Logger
}

func New(graph clienter, publisher event.Publisher, logger logging.Logger) *Service {
	return &Service{
		graph:     graph,
		publisher: publisher,
		logger:    logger,
	}
}

// normalizePath ensures paths like "Documents", "/Documents", "Documents/" all become "/Documents"
func normalizePath(p string) string {
	if p == "" || p == "/" || p == "." {
		return ""
	}
	p = path.Clean("/" + p)
	return p
}

// getUserDriveID retrieves the user's default drive ID
func (s *Service) getUserDriveID(ctx context.Context) (string, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return "", errors.Join(errors.New("unable to instantiate client"), err)
	}

	drive, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		s.logger.Error("failed to retrieve user drive", logging.Any("error", err))
		return "", fmt.Errorf("failed to retrieve user drive: %w", err)
	}

	id := drive.GetId()
	if id == nil || *id == "" {
		s.logger.Error("user drive ID is empty")
		return "", fmt.Errorf("user drive ID is empty")
	}

	s.logger.Info("retrieved user drive ID")
	s.logger.Debug("drive_id", logging.String("id", *id))

	return *id, nil
}

func (s *Service) getChildren(ctx context.Context, folderPath string) (models.DriveItemCollectionResponseable, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return nil, err
	}

	normalized := normalizePath(folderPath)
	s.logger.Debug("normalized folder path", logging.String("path", normalized))

	driveID, err := s.getUserDriveID(ctx)
	if err != nil {
		return nil, err
	}

	var rawURL string
	if normalized != "" {
		rawURL = fmt.Sprintf(rootRelativeChildrenURITemplate, driveID, normalized)
	} else {
		rawURL = fmt.Sprintf(rootChildrenURITemplate, driveID)
	}

	s.logger.Debug("constructed children request URL", logging.String("url", rawURL))

	resp, err := drives.
		NewItemItemsRequestBuilder(rawURL, client.RequestAdapter).
		Get(ctx, nil)

	if err != nil {
		s.logger.Error("failed to retrieve children", logging.String("path", normalized), logging.Any("error", err))
		return nil, fmt.Errorf("failed to retrieve children: %w", err)
	}

	// Extract items for event publishing
	items := resp.GetValue()

	if s.publisher == nil {
		s.logger.Warn("no event publisher configured; skipping drive.children.loaded event")
	} else {
		s.logger.Debug("publishing drive.children.loaded event")
		if err := s.publisher.Publish(ctx, newDriveChildrenLoadedEvent(normalized, items)); err != nil {
			s.logger.Error("failed to publish drive.children.loaded event", logging.Any("error", err))
		}
	}

	s.logger.Info("retrieved drive children successfully", logging.String("path", normalized))
	s.logger.Debug("children_count", logging.Any("count", len(items)))

	return resp, nil
}

func (s *Service) ChildrenIterator(ctx context.Context, folderPath string) iter.Seq2[models.DriveItemable, error] {
	resp, err := s.getChildren(ctx, folderPath)
	if err != nil {
		s.logger.Error("unable to retrieve children", logging.String("path", folderPath), logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, fmt.Errorf("unable to retrieve children: %w", err))
		}
	}

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("unable to instantiate graph client", logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, err)
		}
	}

	pageIterator, err := msgraphcore.NewPageIterator[models.DriveItemable](
		resp,
		client.GetAdapter(),
		models.CreateDriveItemFromDiscriminatorValue,
	)
	if err != nil {
		s.logger.Error("unable to create page iterator", logging.Any("error", err))
		return func(yield func(models.DriveItemable, error) bool) {
			yield(nil, fmt.Errorf("unable to create page iterator: %w", err))
		}
	}

	s.logger.Info("iterating drive children", logging.String("path", folderPath))

	return func(yield func(models.DriveItemable, error) bool) {
		iterErr := pageIterator.Iterate(ctx, func(item models.DriveItemable) bool {
			return yield(item, nil)
		})

		if iterErr != nil {
			s.logger.Error("error during children iteration", logging.Any("error", iterErr))
			yield(nil, iterErr)
		}
	}
}
