package app

import (
	"context"

	domaindrive "github.com/michaeldcanady/go-onedrive/internal/drive/domain"
	domainstate "github.com/michaeldcanady/go-onedrive/internal/state/domain"
)

var _ domaindrive.DriveAliasService = (*AliasService)(nil)

type AliasService struct {
	stateSvc domainstate.Service
}

func NewAliasService(stateSvc domainstate.Service) *AliasService {
	return &AliasService{
		stateSvc: stateSvc,
	}
}

// Resolve resolves the alias to its drive id.
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
