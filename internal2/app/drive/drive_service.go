package drive

import (
	"context"
	"errors"
	"strings"

	domaingraph "github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type driveService struct {
	graph  domaingraph.ClientProvider
	logger logging.Logger
}

func NewDriveService(graph domaingraph.ClientProvider, log logging.Logger) *driveService {
	return &driveService{graph: graph, logger: log}
}

const (
	eventDriveListStart   = "drive.list.start"
	eventDriveListSuccess = "drive.list.success"
	eventDriveListFailure = "drive.list.failure"

	eventDriveResolveStart    = "drive.resolve.start"
	eventDriveResolveMatch    = "drive.resolve.match"
	eventDriveResolveNotFound = "drive.resolve.not_found"
	eventDriveResolveFailure  = "drive.resolve.failure"

	eventDrivePersonalStart   = "drive.personal.start"
	eventDrivePersonalSuccess = "drive.personal.success"
	eventDrivePersonalFailure = "drive.personal.failure"
)

func (s *driveService) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)

	logger.Info("listing drives",
		logging.String("event", eventDriveListStart),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		logger.Error("failed to create graph client",
			logging.String("event", eventDriveListFailure),
			logging.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drives().Get(ctx, nil)
	if err != nil {
		logger.Error("failed to retrieve drives",
			logging.String("event", eventDriveListFailure),
			logging.Error(err),
		)
		return nil, mapGraphError(err)
	}

	out := make([]*drive.Drive, 0, len(resp.GetValue()))
	for _, d := range resp.GetValue() {
		out = append(out, toDomainDrive(d))
	}

	logger.Info("drive list retrieved successfully",
		logging.String("event", eventDriveListSuccess),
		logging.Int("count", len(out)),
	)

	return out, nil
}

func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
		logging.String("drive_ref", driveRef),
	)

	logger.Info("resolving drive reference",
		logging.String("event", eventDriveResolveStart),
	)

	drives, err := s.ListDrives(ctx)
	if err != nil {
		logger.Error("failed to list drives",
			logging.String("event", eventDriveResolveFailure),
			logging.Error(err),
		)
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			logger.Info("drive reference resolved",
				logging.String("event", eventDriveResolveMatch),
				logging.String("drive_id", d.ID),
				logging.String("drive_name", d.Name),
			)
			return d, nil
		}
	}

	logger.Warn("drive reference not found",
		logging.String("event", eventDriveResolveNotFound),
	)

	return nil, errors.New("not found")
}

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	correlationID := util.CorrelationIDFromContext(ctx)

	logger := s.logger.WithContext(ctx).With(
		logging.String("correlation_id", correlationID),
	)

	logger.Info("resolving personal drive",
		logging.String("event", eventDrivePersonalStart),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		logger.Error("failed to create graph client",
			logging.String("event", eventDrivePersonalFailure),
			logging.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		logger.Error("failed to retrieve personal drive",
			logging.String("event", eventDrivePersonalFailure),
			logging.Error(err),
		)
		return nil, mapGraphError(err)
	}

	d := toDomainDrive(resp)

	logger.Info("personal drive resolved successfully",
		logging.String("event", eventDrivePersonalSuccess),
		logging.String("drive_id", d.ID),
		logging.String("drive_name", d.Name),
	)

	return d, nil
}
