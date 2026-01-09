package fsstore

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type FSStore struct {
	baseDir string
	// optional: permissions
	dirPerm  fs.FileMode
	filePerm fs.FileMode
}

func New(baseDir string) *FSStore {
	return &FSStore{
		baseDir:  baseDir,
		dirPerm:  0o700,
		filePerm: 0o600,
	}
}

func (s *FSStore) ensureDir() error {
	return os.MkdirAll(s.baseDir, s.dirPerm)
}

func (s *FSStore) filePathForKey(key string) string {
	// no need to be fancy yet; you can sanitize later
	return filepath.Join(s.baseDir, key)
}

func (s *FSStore) LoadBytes(ctx context.Context, key string) ([]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if err := s.ensureDir(); err != nil {
		return nil, fmt.Errorf("ensure dir: %w", err)
	}

	path := s.filePathForKey(key)
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// treat as cache miss
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read file %s: %w", path, err)
	}
	return data, nil
}

func (s *FSStore) SaveBytes(ctx context.Context, key string, data []byte) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.ensureDir(); err != nil {
		return fmt.Errorf("ensure dir: %w", err)
	}

	path := s.filePathForKey(key)

	// write atomically
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, s.filePerm); err != nil {
		return fmt.Errorf("write tmp file %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("rename tmp file %s -> %s: %w", tmpPath, path, err)
	}

	return nil
}

func (s *FSStore) Delete(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := s.ensureDir(); err != nil {
		return fmt.Errorf("unable to ensure dir: %w", err)
	}

	path := s.filePathForKey(key)
	if err := os.Remove(path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			// already gone
			return nil
		}
		return fmt.Errorf("unable to remove file %s: %w", path, err)
	}
	return nil
}
