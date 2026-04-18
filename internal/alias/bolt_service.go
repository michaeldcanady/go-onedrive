package alias

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	bolt "go.etcd.io/bbolt"
)

type DefaultService struct {
	repo Repository
	db   *bolt.DB
	log  logger.Logger
}

func NewDefaultService(env environment.Service, log logger.Logger) (*DefaultService, error) {
	dbPath, err := env.StateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	dbFilePath := filepath.Join(dbPath, "alias.db")
	db, err := bolt.Open(dbFilePath, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("drive_aliases"))
		return err
	})
	if err != nil {
		db.Close()
		return nil, err
	}

	repo := NewBoltRepository(db)
	return &DefaultService{
		repo: repo,
		db:   db,
		log:  log,
	}, nil
}

func (s *DefaultService) GetDriveIDByAlias(ctx context.Context, name string) (string, error) {
	return s.repo.Get(ctx, name)
}

func (s *DefaultService) GetAliasByDriveID(ctx context.Context, driveID string) (string, error) {
	aliases, err := s.repo.List(ctx)
	if err != nil {
		return "", err
	}
	for al, id := range aliases {
		if id == driveID {
			return al, nil
		}
	}
	return "", fmt.Errorf("alias not found for drive ID %s", driveID)
}

func (s *DefaultService) SetAlias(ctx context.Context, name, driveID string) error {
	return s.repo.Set(ctx, name, driveID)
}

func (s *DefaultService) DeleteAlias(ctx context.Context, name string) error {
	return s.repo.Delete(ctx, name)
}

func (s *DefaultService) ListAliases(ctx context.Context) (map[string]string, error) {
	return s.repo.List(ctx)
}

func (s *DefaultService) Close() error {
	return s.db.Close()
}
