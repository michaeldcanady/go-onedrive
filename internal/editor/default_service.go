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

	"github.com/google/shlex"
	"github.com/michaeldcanady/go-onedrive/internal/environment"
	"github.com/michaeldcanady/go-onedrive/internal/logger"
)

// Option represents a configuration option for the editor service.
type Option func(*DefaultService)

// WithIO sets the standard input, output, and error writers for the editor.
func WithIO(stdin io.Reader, stdout, stderr io.Writer) Option {
	return func(s *DefaultService) {
		s.stdin = stdin
		s.stdout = stdout
		s.stderr = stderr
	}
}

// WithEditor sets an explicit editor command to use.
func WithEditor(cmd string) Option {
	return func(s *DefaultService) {
		s.editorCmd = cmd
	}
}

// DefaultService provides the default implementation of the editor service.
type DefaultService struct {
	envSvc    environment.Service
	log       logger.Logger
	stdin     io.Reader
	stdout    io.Writer
	stderr    io.Writer
	editorCmd string
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(envSvc environment.Service, l logger.Logger, opts ...Option) *DefaultService {
	s := &DefaultService{
		envSvc: envSvc,
		log:    l,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithOptions returns a new Service instance with the specified options applied.
func (s *DefaultService) WithOptions(opts ...Option) Service {
	newS := &DefaultService{
		envSvc:    s.envSvc,
		log:       s.log,
		stdin:     s.stdin,
		stdout:    s.stdout,
		stderr:    s.stderr,
		editorCmd: s.editorCmd,
	}

	for _, opt := range opts {
		opt(newS)
	}

	return newS
}

// Launch launches an external editor with the specified path.
func (s *DefaultService) Launch(ctx context.Context, path string) error {
	editorParts, err := s.getEditorParts()
	if err != nil {
		return err
	}

	args := append(editorParts, path)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	return cmd.Run()
}

// LaunchTempFile creates a temporary file, writes content, launches editor, and returns modified content.
func (s *DefaultService) LaunchTempFile(ctx context.Context, prefix, suffix string, r io.Reader) ([]byte, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read data: %w", err)
	}

	tempDir, _ := s.envSvc.TempDir()
	tmpFile, err := os.CreateTemp(tempDir, prefix+"-*"+suffix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer func() {
		_ = os.Remove(tmpPath)
	}()

	if _, err := tmpFile.Write(data); err != nil {
		_ = tmpFile.Close()
		return nil, fmt.Errorf("failed to write to temp file: %w", err)
	}
	if err := tmpFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	initialHash := sha256.Sum256(data)

	if err := s.Launch(ctx, tmpPath); err != nil {
		return nil, fmt.Errorf("failed to launch editor: %w", err)
	}

	newData, err := os.ReadFile(tmpPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read modified file: %w", err)
	}

	newHash := sha256.Sum256(newData)

	if bytes.Equal(initialHash[:], newHash[:]) {
		return nil, nil
	}

	return newData, nil
}

func (s *DefaultService) getEditorParts() ([]string, error) {
	editorCmd, err := s.getEditorCmd()
	if err != nil {
		return nil, err
	}

	parts, err := shlex.Split(editorCmd)
	if err != nil {
		return nil, fmt.Errorf("failed to parse editor command: %w", err)
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("empty editor command detected")
	}

	return parts, nil
}

func (s *DefaultService) getEditorCmd() (string, error) {
	// 1. Explicitly set command
	if strings.TrimSpace(s.editorCmd) != "" {
		return s.editorCmd, nil
	}

	// 2. Try VISUAL
	if visual, err := s.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	// 3. Try EDITOR
	if editor, err := s.envSvc.Editor(); err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 4. System-specific primary defaults
	if s.envSvc.IsWindows() {
		return "notepad.exe", nil
	}

	// 5. Common Terminal Editors
	if s.envSvc.IsLinux() || s.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 6. OS Opener Defaults
	if s.envSvc.IsMac() {
		if path, err := exec.LookPath("open"); err == nil {
			return path + " -W -t", nil
		}
	}
	if s.envSvc.IsLinux() {
		if path, err := exec.LookPath("xdg-open"); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not detect a suitable editor")
}
