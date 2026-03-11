package app

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/cli/util"
	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
)

type driveService struct {
	gateway domaindrive.DriveGateway
	log     domainlogger.Logger
}

func NewDriveService(gateway domaindrive.DriveGateway, l domainlogger.Logger) *driveService {
	return &driveService{gateway: gateway, log: l}
}

const (
	eventDriveListStart   = "domain.list.start"
	eventDriveListSuccess = "domain.list.success"
	eventDriveListFailure = "domain.list.failure"

	eventDriveResolveStart    = "domain.resolve.match"
	eventDriveResolveMatch    = "domain.resolve.match"
	eventDriveResolveNotFound = "domain.resolve.not_found"
	eventDriveResolveFailure  = "domain.resolve.failure"

	eventDrivePersonalStart   = "domain.personal.start"
	eventDrivePersonalSuccess = "domain.personal.success"
	eventDrivePersonalFailure = "domain.personal.failure"
)

func (s *driveService) ListDrives(ctx context.Context) ([]*domaindrive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		domainlogger.String("correlation_id", correlationID),
	)

	log.Info("listing drives",
		domainlogger.String("event", eventDriveListStart),
	)

	out, err := s.gateway.ListDrives(ctx)
	if err != nil {
		log.Error("failed to retrieve drives",
			domainlogger.String("event", eventDriveListFailure),
			domainlogger.Error(err),
		)
		return nil, err
	}

	log.Info("drive list retrieved successfully",
		domainlogger.String("event", eventDriveListSuccess),
		domainlogger.Int("count", len(out)),
	)

	return out, nil
}

func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*domaindrive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		domainlogger.String("correlation_id", correlationID),
		domainlogger.String("drive_ref", driveRef),
	)

	log.Info("resolving drive reference",
		domainlogger.String("event", eventDriveResolveStart),
	)

	drives, err := s.ListDrives(ctx)
	if err != nil {
		log.Error("failed to list drives",
			domainlogger.String("event", eventDriveResolveFailure),
			domainlogger.Error(err),
		)
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			log.Info("drive reference resolved",
				domainlogger.String("event", eventDriveResolveMatch),
				domainlogger.String("drive_id", d.ID),
				domainlogger.String("drive_name", d.Name),
			)
			return d, nil
		}
	}

	log.Warn("drive reference not found",
		domainlogger.String("event", eventDriveResolveNotFound),
	)

	return nil, errors.New("not found")
}

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*domaindrive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	log := s.log.WithContext(ctx).With(
		domainlogger.String("correlation_id", correlationID),
	)

	log.Info("resolving personal drive",
		domainlogger.String("event", eventDrivePersonalStart),
	)

	d, err := s.gateway.GetPersonalDrive(ctx)
	if err != nil {
		log.Error("failed to retrieve personal drive",
			domainlogger.String("event", eventDrivePersonalFailure),
			domainlogger.Error(err),
		)
		return nil, err
	}

	log.Info("personal drive resolved successfully",
		domainlogger.String("event", eventDrivePersonalSuccess),
		domainlogger.String("drive_id", d.ID),
		domainlogger.String("drive_name", d.Name),
	)

	return d, nil
}
