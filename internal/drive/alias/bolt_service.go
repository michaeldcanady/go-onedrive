package alias

import (
	"errors"
	"path/filepath"
	"time"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	coreerrors "github.com/michaeldcanady/go-onedrive/internal/errors"
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
		return nil, coreerrors.NewInternal(err, "failed to get state directory", "Ensure the application has proper permissions to access its state directory.")
	}

	dbFilePath := filepath.Join(dbPath, aliasDBFileName)
	db, err := bbolt.Open(dbFilePath, 0600, &bbolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, coreerrors.NewInternal(err, "failed to open drive alias database", "Check if another instance of the application is running or if the state directory is accessible.")
	}

	bs := &BoltService{
		db:  db,
		log: log,
	}

	if err := bs.ensureBucket(); err != nil {
		return nil, err
	}

	return bs, nil
}

// ensureBucket creates the necessary bucket for storing drive aliases if it does not already exist.
func (s *BoltService) ensureBucket() error {
	return s.db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(driveAliasesBucketName)
		if err != nil {
			return coreerrors.NewWriteError(err, "failed to create drive aliases bucket", "This may indicate a file system error or disk space issue.")
		}
		return nil
	})
}

func (s *BoltService) getBucket(tx *bbolt.Tx) (*bbolt.Bucket, error) {
	b := tx.Bucket(driveAliasesBucketName)
	if b == nil {
		s.log.Error("drive aliases bucket not found", logger.String("bucket", string(driveAliasesBucketName)))
		return nil, NewBucketNotFoundError()
	}
	return b, nil
}

// Close closes the BoltDB database connection.
func (s *BoltService) Close() error {
	s.log.Debug("closing drive alias database")
	return s.db.Close()
}

// GetAliasByDriveID retrieves the alias for a given drive ID, if it exists.
func (s *BoltService) GetAliasByDriveID(input string) (string, error) {
	s.log.Debug("resolving alias by drive ID", logger.String("drive_id", input))
	var alias string
	err := s.iterateAliases(func(al, driveID string) error {
		if driveID == input {
			alias = al
			return nil // Stop iteration once we find the drive ID
		}
		return nil
	})
	if err != nil {
		s.log.Error("failed to iterate aliases by drive ID", logger.Error(err))
	}
	return alias, err
}

// GetDriveIDByAlias retrieves the drive ID associated with a given alias, if it exists.
func (s *BoltService) GetDriveIDByAlias(alias string) (string, error) {
	s.log.Debug("resolving drive ID by alias", logger.String("alias", alias))
	var driveID string
	err := s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := s.getBucket(tx)
		if err != nil {
			return err
		}
		v := bucket.Get([]byte(alias))
		if v != nil {
			driveID = string(v)
		}
		return nil
	})
	if err != nil {
		if !errors.Is(err, coreerrors.CodeNotFound) && !errors.Is(err, coreerrors.CodeInternal) {
			s.log.Error("failed to resolve drive ID by alias", logger.String("alias", alias), logger.Error(err))
		}
		return "", err
	}
	if driveID == "" {
		s.log.Debug("alias not found", logger.String("alias", alias))
		return "", NewAliasNotFoundError(alias)
	}
	return driveID, nil
}

// SetAlias assigns an alias to a specific drive ID in the BoltDB.
func (s *BoltService) SetAlias(driveID string, alias string) error {
	s.log.Info("setting drive alias", logger.String("drive_id", driveID), logger.String("alias", alias))
	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := s.getBucket(tx)
		if err != nil {
			return err
		}
		err = bucket.Put([]byte(alias), []byte(driveID))
		if err != nil {
			return coreerrors.NewWriteError(err, "failed to save drive alias", "Try again or check if the database file is read-only.")
		}
		return nil
	})
}

// DeleteAlias removes the alias associated with a specific drive ID from the BoltDB.
func (s *BoltService) DeleteAlias(alias string) error {
	s.log.Info("deleting drive alias", logger.String("alias", alias))

	return s.db.Update(func(tx *bbolt.Tx) error {
		bucket, err := s.getBucket(tx)
		if err != nil {
			return err
		}
		err = bucket.Delete([]byte(alias))
		if err != nil {
			return coreerrors.NewWriteError(err, "failed to delete drive alias", "Try again or check if the database file is read-only.")
		}
		return nil
	})
}

// ListAliases retrieves all drive IDs and their corresponding aliases from the BoltDB.
func (s *BoltService) ListAliases() (map[string]string, error) {
	s.log.Debug("listing drive aliases")
	aliases := make(map[string]string)
	err := s.iterateAliases(func(alias string, driveID string) error {
		aliases[alias] = driveID
		return nil
	})
	if err != nil {
		s.log.Error("failed to list drive aliases", logger.Error(err))
	}
	return aliases, err
}

func (s *BoltService) iterateAliases(fn func(alias string, driveID string) error) error {
	return s.db.View(func(tx *bbolt.Tx) error {
		bucket, err := s.getBucket(tx)
		if err != nil {
			return err
		}
		return bucket.ForEach(func(k, v []byte) error {
			return fn(string(k), string(v))
		})
	})
}
