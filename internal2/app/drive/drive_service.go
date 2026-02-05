package drive

import (
	"context"
	"errors"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

type driveService struct {
	graph  clienter
	logger logging.Logger
}

func NewDriveService(graph clienter, log logging.Logger) *driveService {
	return &driveService{graph: graph, logger: log}
}

func (s *driveService) ListDrives(ctx context.Context) ([]*drive.Drive, error) {
	s.logger.Debug("listing drives")

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("failed to create graph client",
			logging.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drives().Get(ctx, nil)
	if err != nil {
		mapped := mapGraphError(err)
		s.logger.Error("graph error while listing drives",
			logging.Error(mapped),
		)
		return nil, mapped
	}

	values := resp.GetValue()
	s.logger.Info("drives retrieved successfully",
		logging.Int("count", len(values)),
	)

	out := make([]*drive.Drive, 0, len(values))
	for _, d := range values {
		out = append(out, toDomainDrive(d))
	}

	return out, nil
}

func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*drive.Drive, error) {
	s.logger.Debug("resolving drive",
		logging.String("driveRef", driveRef),
	)

	drives, err := s.ListDrives(ctx)
	if err != nil {
		s.logger.Error("failed to list drives during resolve",
			logging.String("driveRef", driveRef),
			logging.Error(err),
		)
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			s.logger.Info("drive resolved successfully",
				logging.String("driveRef", driveRef),
				logging.String("resolvedID", d.ID),
				logging.String("resolvedName", d.Name),
			)
			return d, nil
		}
	}

	s.logger.Warn("drive not found",
		logging.String("driveRef", driveRef),
	)

	return nil, errors.New("not found")
}

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	s.logger.Debug("resolving personal drive")

	client, err := s.graph.Client(ctx)
	if err != nil {
		s.logger.Error("failed to create graph client",
			logging.Error(err),
		)
		return nil, err
	}

	resp, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		mapped := mapGraphError(err)
		s.logger.Error("graph error while resolving personal drive",
			logging.Error(mapped),
		)
		return nil, mapped
	}

	d := toDomainDrive(resp)

	s.logger.Info("personal drive resolved successfully",
		logging.String("driveID", d.ID),
		logging.String("driveName", d.Name),
	)

	return d, nil
}
