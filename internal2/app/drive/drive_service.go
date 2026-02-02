package drive

import (
	"context"
	"errors"
	"fmt"
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
	client, err := s.graph.Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.Me().Drives().Get(ctx, nil)
	if err != nil {
		return nil, mapGraphError(err)
	}

	out := make([]*drive.Drive, 0, len(resp.GetValue()))
	for _, d := range resp.GetValue() {
		fmt.Println(deref(d.GetName()))
		out = append(out, toDomainDrive(d))
	}

	return out, nil
}

func (s *driveService) ResolveDrive(ctx context.Context, driveRef string) (*drive.Drive, error) {
	drives, err := s.ListDrives(ctx)
	if err != nil {
		return nil, err
	}

	for _, d := range drives {
		if strings.EqualFold(d.ID, driveRef) || strings.EqualFold(d.Name, driveRef) {
			return d, nil
		}
	}

	return nil, errors.New("not found")
}

func (s *driveService) ResolvePersonalDrive(ctx context.Context) (*drive.Drive, error) {
	client, err := s.graph.Client(ctx)
	if err != nil {
		return nil, err
	}

	resp, err := client.Me().Drive().Get(ctx, nil)
	if err != nil {
		return nil, mapGraphError(err)
	}

	return toDomainDrive(resp), nil
}
