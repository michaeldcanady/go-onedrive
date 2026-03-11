package app

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"os"

	domainlogger "github.com/michaeldcanady/go-onedrive/internal/core/logger/domain"
	domaineditor "github.com/michaeldcanady/go-onedrive/internal/editor/domain"
)

var _ domaineditor.Service = (*Service)(nil)

type Service struct {
	launcher domaineditor.Launcher
	log      domainlogger.Logger
	stdin    io.Reader
	stdout   io.Writer
	stderr   io.Writer
}

func NewService(launcher domaineditor.Launcher, l domainlogger.Logger) *Service {
	return &Service{
		launcher: launcher,
		log:      l,
		stdin:    os.Stdin,
		stdout:   os.Stdout,
		stderr:   os.Stderr,
	}
}

func (s *Service) WithIO(stdin io.Reader, stdout, stderr io.Writer) domaineditor.Service {
	return &Service{
		launcher: s.launcher.WithIO(stdin, stdout, stderr),
		log:      s.log,
		stdin:    stdin,
		stdout:   stdout,
		stderr:   stderr,
	}
}

func (s *Service) Launch(path string) error {
	return s.launcher.Launch(path)
}

func (s *Service) LaunchTempFile(prefix, suffix string, r io.Reader) ([]byte, string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read data: %w", err)
	}

	tmpFile, err := os.CreateTemp("", prefix+"-*"+suffix)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return nil, "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, "", fmt.Errorf("failed to close temp file: %w", err)
	}

	initialHash := sha256.Sum256(data)

	if err := s.Launch(tmpPath); err != nil {
		return nil, "", fmt.Errorf("failed to launch editor: %w", err)
	}

	newData, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read modified file: %w", err)
	}

	newHash := sha256.Sum256(newData)

	if bytes.Equal(initialHash[:], newHash[:]) {
		return nil, "", nil
	}

	return newData, tmpPath, nil
}
