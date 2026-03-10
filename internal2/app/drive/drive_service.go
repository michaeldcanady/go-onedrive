package drive

import (
	"context"
	"errors"
	"strings"

	domaingraph "github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type driveService struct {
	graph domaingraph.ClientProvider
	log   logger.Logger
}

func NewDriveService(graph domaingraph.ClientProvider, l logger.Logger) *driveService {
	return &driveService{graph: graph, log: l}
}

const (
	eventDriveListStart   = "drive.list.start"
	eventDriveListSuccess = "drive.list.success"
	eventDriveListFailure = "drive.list.failure"

	eventDriveResolveStart    = "drive.resolve.match"
	eventDriveResolveMatch    = "drive.resolve.match"
	eventDriveResolveNotFound = "drive.resolve.not_found"
	eventDriveResolveFailure  = "drive.resolve.failure"

	eventDrivePersonalStart   = "drive.personal.start"
	eventDrivePersonalSuccess = "drive.personal.success"
	eventDrivePersonalFailure = "drive.personal.failure"
)

// ListDrives lists available onedrive drives.
func (s *driveService) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)

	log.Info("listing drives",
		logger.String("event", eventDriveListStart),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		log.Error("failed to create graph client",
			logger.String("event", eventDriveListFailure),
			logger.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drives().Get(ctx, nil)
	if err != nil {
		log.Error("failed to retrieve drives",
			logger.String("event", eventDriveListFailure),
			logger.Error(err),
		)
		return nil, mapGraphError(err)
	}

	out := make([]*drive.Drive, 0, len(resp.GetValue()))
	for _, d := range resp.GetValue() {
		out = append(out, toDomainDrive(d))
	}

	log.Info("drive list retrieved successfully",
		logger.String("event", eventDriveListSuccess),
		logger.Int("count", len(out)),
	)

	return out, nil
}

// ResolveDrive resolves drive from ref (id).
func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
		logger.String("drive_ref", driveRef),
	)

	log.Info("resolving drive reference",
		logger.String("event", eventDriveResolveStart),
	)

	drives, err := s.ListDrives(ctx)
	if err != nil {
		log.Error("failed to list drives",
			logger.String("event", eventDriveResolveFailure),
			logger.Error(err),
		)
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			log.Info("drive reference resolved",
				logger.String("event", eventDriveResolveMatch),
				logger.String("drive_id", d.ID),
				logger.String("drive_name", d.Name),
			)
			return d, nil
		}
	}

	log.Warn("drive reference not found",
		logger.String("event", eventDriveResolveNotFound),
	)

	return nil, errors.New("not found")
}

// ResolvePersonalDrive resolves logged in user
func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)

	log.Info("resolving personal drive",
		logger.String("event", eventDrivePersonalStart),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		log.Error("failed to create graph client",
			logger.String("event", eventDrivePersonalFailure),
			logger.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		log.Error("failed to retrieve personal drive",
			logger.String("event", eventDrivePersonalFailure),
			logger.Error(err),
		)
		return nil, mapGraphError(err)
	}

	d := toDomainDrive(resp)

	log.Info("personal drive resolved successfully",
		logger.String("event", eventDrivePersonalSuccess),
		logger.String("drive_id", d.ID),
		logger.String("drive_name", d.Name),
	)

	return d, nil
}
