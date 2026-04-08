package editor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

const (
	defaultUnixShell    = "/bin/bash"
	defaultWindowsShell = "cmd"
)

// DefaultService provides the default implementation of the editor service.
type DefaultService struct {
	envSvc environment.Service
	log    logger.Logger
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(envSvc environment.Service, l logger.Logger) *DefaultService {
	return &DefaultService{
		envSvc: envSvc,
		log:    l,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// WithIO returns a new Service instance with the specified standard input, output, and error writers.
func (s *DefaultService) WithIO(stdin io.Reader, stdout, stderr io.Writer) Service {
	return &DefaultService{
		envSvc: s.envSvc,
		log:    s.log,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}
}

// Launch launches an external editor with the specified path.
func (s *DefaultService) Launch(path string) error {
	s.log.Info("launching external editor", logger.String("path", path))

	args, err := s.buildCmd(path)
	if err != nil {
		s.log.Error("failed to build editor command", logger.Error(err))
		return err
	}

	s.log.Debug("executing editor command", logger.String("command", strings.Join(args, " ")))

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	if err := cmd.Run(); err != nil {
		s.log.Error("editor process exited with error", logger.Error(err), logger.String("path", path))
		return err
	}

	s.log.Debug("editor process finished successfully", logger.String("path", path))
	return nil
}

// LaunchTempFile creates a temporary file, writes content, launches editor, and returns modified content.
func (s *DefaultService) LaunchTempFile(ctx context.Context, prefix, suffix string, r io.Reader) ([]byte, string, error) {
	l := s.log.WithContext(ctx)
	l.Debug("preparing temporary file for editing", logger.String("prefix", prefix), logger.String("suffix", suffix))

	data, err := io.ReadAll(r)
	if err != nil {
		l.Error("failed to read data for temp file", logger.Error(err))
		return nil, "", fmt.Errorf("failed to read data: %w", err)
	}

	tempDir, err := s.envSvc.TempDir()
	if err != nil {
		l.Error("failed to get temp directory", logger.Error(err))
		return nil, "", fmt.Errorf("failed to get temp dir: %w", err)
	}

	tmpFile, err := os.CreateTemp(tempDir, prefix+"-*"+suffix)
	if err != nil {
		l.Error("failed to create temp file", logger.Error(err))
		return nil, "", fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	l.Debug("temporary file created", logger.String("path", tmpPath))
	defer func() {
		l.Debug("cleaning up temporary file", logger.String("path", tmpPath))
		if err := os.Remove(tmpPath); err != nil {
			l.Warn("failed to remove temp file", logger.String("path", tmpPath), logger.Error(err))
		}
	}()

	if _, err := tmpFile.Write(data); err != nil {
		l.Error("failed to write to temp file", logger.String("path", tmpPath), logger.Error(err))
		_ = tmpFile.Close()
		return nil, "", fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		l.Error("failed to close temp file handle", logger.String("path", tmpPath), logger.Error(err))
		return nil, "", fmt.Errorf("failed to close temp file: %w", err)
	}

	initialHash := sha256.Sum256(data)

	if err := s.Launch(tmpPath); err != nil {
		return nil, "", fmt.Errorf("failed to launch editor: %w", err)
	}

	newData, err := os.ReadFile(tmpPath)
	if err != nil {
		l.Error("failed to read modified file", logger.String("path", tmpPath), logger.Error(err))
		return nil, "", fmt.Errorf("failed to read modified file: %w", err)
	}

	newHash := sha256.Sum256(newData)

	if bytes.Equal(initialHash[:], newHash[:]) {
		l.Info("no modifications detected", logger.String("path", tmpPath))
		return nil, "", nil
	}

	l.Info("file modifications detected", logger.String("path", tmpPath))
	return newData, tmpPath, nil
}

// buildCmd constructs the command to launch the editor based on environment variables and OS defaults.
func (s *DefaultService) buildCmd(path string) ([]string, error) {
	shell := s.getShell()

	editorCmd, err := s.getEditor()
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s %s", editorCmd, path)
	return append(shell, cmd), nil
}

// getEditor determines the appropriate editor command based on environment variables and OS defaults.
func (s *DefaultService) getEditor() (string, error) {
	// 1. Try VISUAL
	if visual, err := s.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		s.log.Debug("using editor from VISUAL environment variable", logger.String("command", visual))
		return visual, nil
	}

	// 2. Try EDITOR
	editor, err := s.envSvc.Editor()
	if err == nil && strings.TrimSpace(editor) != "" {
		s.log.Debug("using editor from EDITOR environment variable", logger.String("command", editor))
		return editor, nil
	}

	// 3. System-specific primary defaults
	if s.envSvc.IsWindows() {
		s.log.Debug("using default Windows editor", logger.String("command", "notepad"))
		return "notepad", nil
	}

	// 4. Common Terminal Editors
	if s.envSvc.IsLinux() || s.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				s.log.Debug("using fallback terminal editor", logger.String("command", path))
				return path, nil
			}
		}
	}

	// 5. OS Opener Defaults
	if s.envSvc.IsMac() {
		if path, err := exec.LookPath("open"); err == nil {
			cmd := path + " -W -t"
			s.log.Debug("using macOS 'open' as editor", logger.String("command", cmd))
			return cmd, nil
		}
	}
	if s.envSvc.IsLinux() {
		if path, err := exec.LookPath("xdg-open"); err == nil {
			s.log.Debug("using Linux 'xdg-open' as editor", logger.String("command", path))
			return path, nil
		}
	}

	s.log.Error("failed to detect suitable editor")
	return "", fmt.Errorf("could not detect a suitable editor")
}

// getShell returns the appropriate shell command based on the operating system and environment variables.
func (s *DefaultService) getShell() []string {
	if s.envSvc.IsWindows() {
		return []string{defaultWindowsShell, "/C"}
	}

	shell, err := s.envSvc.Shell()
	if err != nil || strings.TrimSpace(shell) == "" {
		shell = defaultUnixShell
	}

	return []string{shell, "-c"}
}
