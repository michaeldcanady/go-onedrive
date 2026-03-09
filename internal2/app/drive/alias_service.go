package drive

import (
	"context"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/drive"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/state"
)

var _ drive.DriveAliasService = (*AliasService)(nil)

type AliasService struct {
	stateSvc state.Service
}

func NewAliasService(stateSvc state.Service) *AliasService {
	return &AliasService{
		stateSvc: stateSvc,
	}
}

func (s *AliasService) Resolve(ctx context.Context, alias string) (string, error) {
	// 1. Try to get driveID from state (persistent aliases)
	driveID, err := s.stateSvc.GetDriveAlias(alias)
	if err != nil {
		return "", err
	}

	if driveID != "" {
		return driveID, nil
	}

	// 2. If not an alias, assume it's a Drive ID itself
	// (Alternatively, we could verify it against Graph, but for now we trust it)
	return alias, nil
}
