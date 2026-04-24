package editor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// CreateSession initializes a new editing session.
func (s *DefaultService) CreateSession(ctx context.Context, remoteURI *fs.URI, r io.Reader) (*Session, error) {
	tempDir, err := s.envSvc.TempDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get temp directory: %w", err)
	}

	ext := filepath.Ext(remoteURI.Path)
	tmpFile, err := os.CreateTemp(tempDir, "odc-edit-*"+ext)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tmpFile.Close()

	localPath := tmpFile.Name()
	localURI, err := s.uriFactory.FromLocalPath(localPath)
	if err != nil {
		_ = os.Remove(localPath)
		return nil, fmt.Errorf("failed to create local URI: %w", err)
	}

	// Stream and Hash
	hash := sha256.New()
	mw := io.MultiWriter(tmpFile, hash)

	if _, err := io.Copy(mw, r); err != nil {
		_ = os.Remove(localPath)
		return nil, fmt.Errorf("failed to stage content to local file: %w", err)
	}

	session := &Session{
		ID:          uuid.New().String(),
		RemoteURI:   remoteURI,
		LocalURI:    localURI,
		InitialHash: hash.Sum(nil),
		state:       StateCreated,
	}

	return session, nil
}

// Modified checks if the local file in the session has changed.
func (s *DefaultService) Modified(session *Session) (bool, error) {
	if state := session.State(); state != StateCompleted {
		return false, fmt.Errorf("cannot check modifications for session in state %s", state)
	}

	f, err := os.Open(session.LocalURI.Path)
	if err != nil {
		return false, fmt.Errorf("failed to open local file for modification check: %w", err)
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return false, fmt.Errorf("failed to hash local file: %w", err)
	}

	return !bytes.Equal(session.InitialHash, hash.Sum(nil)), nil
}

// NewContent returns a reader for the modified content in the session.
func (s *DefaultService) NewContent(session *Session) (io.ReadCloser, error) {
	if state := session.State(); state != StateCompleted {
		return nil, fmt.Errorf("cannot get content for session in state %s", state)
	}

	f, err := os.Open(session.LocalURI.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open local file: %w", err)
	}
	return f, nil
}

// Cleanup removes the temporary local file and releases session resources.
func (s *DefaultService) Cleanup(ctx context.Context, session *Session) error {
	return session.Handle(ctx, s, EventClose)
}

// removeFile is the internal implementation that actually deletes the local file.
func (s *DefaultService) removeFile(session *Session) error {
	return os.Remove(session.LocalURI.Path)
}
