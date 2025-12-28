package cache

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/gofrs/flock"
)

type FileCache struct {
	path string
}

func NewFileCache(path string) *FileCache {
	return &FileCache{
		path: path,
	}
}

func (c *FileCache) filePath(key string) string {
	return filepath.Join(c.path, key)
}

func (c *FileCache) lockPath(key string) string {
	// Separate lock file to avoid locking the actual data file
	return filepath.Join(c.path, key+".lock")
}

func (c *FileCache) init() error {
	// Ensure directory exists
	if err := os.MkdirAll(c.path, 0o755); err != nil {
		return errors.Join(ErrCreateCacheDir, err)
	}

	return nil
}

// Set writes the content of reader into the cache under the given key.
// Uses an exclusive lock via gofrs/flock.
func (c *FileCache) Set(key string, reader io.Reader) error {
	target := c.filePath(key)
	lockFile := c.lockPath(key)

	if err := c.init(); err != nil {
		return err
	}

	// Acquire exclusive lock
	fl := flock.New(lockFile)
	locked, err := fl.TryLock()
	if err != nil {
		return errors.Join(ErrAcquireLock, err)
	}
	if !locked {
		return ErrAcquireLock
	}
	defer fl.Unlock()

	// Create temp file for atomic write
	tmp, err := os.CreateTemp(filepath.Dir(target), key+".tmp-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	defer os.Remove(tmp.Name())
	defer tmp.Close()

	// Write data
	if _, err := io.Copy(tmp, reader); err != nil {
		return fmt.Errorf("writing temp file: %w", err)
	}

	// Sync to disk
	if err := tmp.Sync(); err != nil {
		return fmt.Errorf("sync temp file: %w", err)
	}

	// Atomic replace
	if err := os.Rename(tmp.Name(), target); err != nil {
		return fmt.Errorf("rename temp file: %w", err)
	}

	return nil
}

// Get reads the cached value for key into writer.
// Uses a shared lock via gofrs/flock.
func (c *FileCache) Get(key string, writer io.Writer) error {
	target := c.filePath(key)
	lockFile := c.lockPath(key)

	// Acquire shared lock
	fl := flock.New(lockFile)
	locked, err := fl.TryRLock()
	if err != nil {
		return errors.Join(ErrAcquireLock, err)
	}
	if !locked {
		return ErrAcquireLock
	}
	defer fl.Unlock()

	// Open file
	f, err := os.Open(target)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache miss: %s", key)
		}
		return fmt.Errorf("opening cache file: %w", err)
	}
	defer f.Close()

	// Copy to writer
	if _, err := io.Copy(writer, f); err != nil {
		return fmt.Errorf("reading cache file: %w", err)
	}

	return nil
}
