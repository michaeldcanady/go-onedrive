package alias

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
	"go.etcd.io/bbolt"
)

const (
	// aliasDBFileName is the name of the BoltDB file used for storing drive alias data.
	aliasDBFileName = "alias.db"
)

var (
	// driveAliasesBucketName is used to store user-defined drive aliases.
	driveAliasesBucketName = []byte("drive_aliases")

	ErrDriveIDNotFound = errors.New("drive ID not found for alias")
)

// BoltService is a persistent implementation of the alias.Service using BoltDB.
type BoltService struct {
	// db is the BoltDB database connection used for storing drive alias data.
	db *bbolt.DB
	// log is the logger used for reporting configuration events.
	log logger.Logger
}

// NewBoltService initializes a new instance of the BoltService with the provided environment and logger.
func NewBoltService(env environment.Service, log logger.Logger) (*BoltService, error) {
	dbPath, err := env.StateDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get state directory: %w", err)
	}

	dbFilePath := filepath.Join(dbPath, aliasDBFileName)
	db, err := bbolt.Open(dbFilePath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("failed to open BoltDB: %w", err)
	}

	bs := &BoltService{
		db:  db,
		log: log,
	}

	if err := bs.ensureBucket(); err != nil {
		return nil, fmt.Errorf("failed to ensure bucket: %w", err)
	}

	return bs, nil
}

// ensureBucket creates the necessary bucket for storing drive aliases if it does not already exist.
func (s *BoltService) ensureBucket() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(driveAliasesBucketName)
		return err
	})
}

// Close closes the BoltDB database connection.
func (s *BoltService) Close() error {
	return s.db.Close()
}

// GetAliasByDriveID retrieves the alias for a given drive ID, if it exists.
func (s *BoltService) GetAliasByDriveID(input string) (string, error) {
	var alias string
	err := s.iterateAliases(func(al, driveID string) error {
		if driveID == input {
			alias = al
			return nil // Stop iteration once we find the drive ID
		}
		return nil
	})
	return alias, err
}

// GetDriveIDByAlias retrieves the drive ID associated with a given alias, if it exists.
func (s *BoltService) GetDriveIDByAlias(alias string) (string, error) {
	var driveID string
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(driveAliasesBucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", driveAliasesBucketName)
		}
		return bucket.ForEach(func(k, v []byte) error {
			if string(v) == alias {
				driveID = string(k)
				return nil // Stop iteration once we find the alias
			}
			return nil
		})
	})
	if err != nil {
		return "", err
	}
	if driveID == "" {
		return "", ErrDriveIDNotFound
	}
	return driveID, nil
}

// SetAlias assigns an alias to a specific drive ID in the BoltDB.
func (s *BoltService) SetAlias(driveID string, alias string) error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(driveAliasesBucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", driveAliasesBucketName)
		}
		return bucket.Put([]byte(alias), []byte(driveID))
	})
}

// DeleteAlias removes the alias associated with a specific drive ID from the BoltDB.
func (s *BoltService) DeleteAlias(alias string) error {
	driveID, err := s.GetDriveIDByAlias(alias)
	if err != nil && !errors.Is(err, ErrDriveIDNotFound) {
		return err
	}

	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(driveAliasesBucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", driveAliasesBucketName)
		}
		return bucket.Delete([]byte(driveID))
	})
}

// ListAliases retrieves all drive IDs and their corresponding aliases from the BoltDB.
func (s *BoltService) ListAliases() (map[string]string, error) {
	aliases := make(map[string]string)
	err := s.iterateAliases(func(alias string, driveID string) error {
		aliases[alias] = driveID
		return nil
	})
	return aliases, err
}

func (s *BoltService) iterateAliases(fn func(alias string, driveID string) error) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket(driveAliasesBucketName)
		if bucket == nil {
			return fmt.Errorf("bucket %s not found", driveAliasesBucketName)
		}
		return bucket.ForEach(func(k, v []byte) error {
			return fn(string(k), string(v))
		})
	})
}
