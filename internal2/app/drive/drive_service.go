package drive

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type driveService struct {
	gateway drive.DriveGateway
	log     logger.Logger
}

func NewDriveService(gateway drive.DriveGateway, l logger.Logger) *driveService {
	return &driveService{gateway: gateway, log: l}
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

func (s *driveService) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)

	log.Info("listing drives",
		logger.String("event", eventDriveListStart),
	)

	out, err := s.gateway.ListDrives(ctx)
	if err != nil {
		log.Error("failed to retrieve drives",
			logger.String("event", eventDriveListFailure),
			logger.Error(err),
		)
		return nil, mapGraphError(err)
	}

	log.Info("drive list retrieved successfully",
		logger.String("event", eventDriveListSuccess),
		logger.Int("count", len(out)),
	)

	return out, nil
}

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

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		logger.String("correlation_id", correlationID),
	)

	log.Info("resolving personal drive",
		logger.String("event", eventDrivePersonalStart),
	)

	d, err := s.gateway.GetPersonalDrive(ctx)
	if err != nil {
		log.Error("failed to retrieve personal drive",
			logger.String("event", eventDrivePersonalFailure),
			logger.Error(err),
		)
		return nil, mapGraphError(err)
	}

	log.Info("personal drive resolved successfully",
		logger.String("event", eventDrivePersonalSuccess),
		logger.String("drive_id", d.ID),
		logger.String("drive_name", d.Name),
	)

	return d, nil
}
