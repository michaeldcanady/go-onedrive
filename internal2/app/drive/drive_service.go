package drive

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/graph"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
	"github.com/michaeldcanady/go-onedrive/internal2/interface/cli/util"
)

type driveService struct {
	graph  graph.ClientProvider
	logger logging.Logger
}

func NewDriveService(graph graph.ClientProvider, log logging.Logger) *driveService {
	return &driveService{graph: graph, logger: log}
}

func (s *driveService) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("listing drives",
		logging.String("event", "drive_list_start"),
		logging.String("correlation_id", cid),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("failed to create graph client",
			logging.String("event", "drive_list_client_error"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	resp, err := client.Me().Drives().Get(ctx, nil)
	if err != nil {
		mapped := mapGraphError(err)
		s.logger.Error("graph error while listing drives",
			logging.String("event", "drive_list_graph_error"),
			logging.Error(mapped),
			logging.String("correlation_id", cid),
		)
		return nil, mapped
	}

	values := resp.GetValue()

	s.logger.Info("drives retrieved successfully",
		logging.String("event", "drive_list_success"),
		logging.Int("count", len(values)),
		logging.String("correlation_id", cid),
	)

	out := make([]*drive.Drive, 0, len(values))
	for _, d := range values {
		out = append(out, toDomainDrive(d))
	}

	return out, nil
}

func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*drive.Drive, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("resolving drive",
		logging.String("event", "drive_resolve_start"),
		logging.String("drive_ref", driveRef),
		logging.String("correlation_id", cid),
	)

	drives, err := s.ListDrives(ctx)
	if err != nil {
		s.logger.Error("failed to list drives during resolve",
			logging.String("event", "drive_resolve_list_error"),
			logging.String("drive_ref", driveRef),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			s.logger.Info("drive resolved successfully",
				logging.String("event", "drive_resolve_success"),
				logging.String("drive_ref", driveRef),
				logging.String("resolved_id", d.ID),
				logging.String("resolved_name", d.Name),
				logging.String("correlation_id", cid),
			)
			return d, nil
		}
	}

	s.logger.Warn("drive not found",
		logging.String("event", "drive_resolve_not_found"),
		logging.String("drive_ref", driveRef),
		logging.String("correlation_id", cid),
	)

	return nil, errors.New("not found")
}

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	cid := util.CorrelationIDFromContext(ctx)

	s.logger.Debug("resolving personal drive",
		logging.String("event", "drive_resolve_personal_start"),
		logging.String("correlation_id", cid),
	)

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("failed to create graph client",
			logging.String("event", "drive_resolve_personal_client_error"),
			logging.Error(err),
			logging.String("correlation_id", cid),
		)
		return nil, err
	}

	resp, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		mapped := mapGraphError(err)
		s.logger.Error("graph error while resolving personal drive",
			logging.String("event", "drive_resolve_personal_graph_error"),
			logging.Error(mapped),
			logging.String("correlation_id", cid),
		)
		return nil, mapped
	}

	d := toDomainDrive(resp)

	s.logger.Info("personal drive resolved successfully",
		logging.String("event", "drive_resolve_personal_success"),
		logging.String("drive_id", d.ID),
		logging.String("drive_name", d.Name),
		logging.String("correlation_id", cid),
	)

	return d, nil
}
