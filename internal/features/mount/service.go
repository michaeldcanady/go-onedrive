package mount

import (
	"context"
	"fmt"

	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
)

type mountService struct {
	repo   Repository
	logger logger.Service
}

// NewMountService returns a new [Service] initialized with the provided repository.
func NewMountService(repo Repository, l logger.Service) Service {
	return &mountService{
		repo:   repo,
		logger: l,
	}
}

func (s *mountService) Add(ctx context.Context, m *Mount) error {
	if err := s.repo.Save(m); err != nil {
		return fmt.Errorf("failed to save mount: %w", err)
	}
	s.logger.Info("mount added", "path", m.Path, "type", m.Type)
	return nil
}

func (s *mountService) List(ctx context.Context) ([]*Mount, error) {
	return s.repo.List()
}

func (s *mountService) Remove(ctx context.Context, path string) error {
	if err := s.repo.Delete(path); err != nil {
		return fmt.Errorf("failed to delete mount: %w", err)
	}
	s.logger.Info("mount removed", "path", path)
	return nil
}

func (s *mountService) Get(ctx context.Context, path string) (*Mount, error) {
	return s.repo.Get(path)
}
