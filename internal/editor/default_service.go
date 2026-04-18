package editor

import (
	"bytes"
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
	args, err := s.buildCmd(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	return cmd.Run()
}

// LaunchTempFile creates a temporary file, writes content, launches editor, and returns modified content.
func (s *DefaultService) LaunchTempFile(prefix, suffix string, r io.Reader) ([]byte, string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read data: %w", err)
	}

	tempDir, _ := s.envSvc.TempDir()
	tmpFile, err := os.CreateTemp(tempDir, prefix+"-*"+suffix)
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

func (s *DefaultService) buildCmd(path string) ([]string, error) {
	shell := s.getShell()

	editorCmd, err := s.getEditor()
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s %s", editorCmd, path)
	return append(shell, cmd), nil
}

func (s *DefaultService) getEditor() (string, error) {
	// 1. Try VISUAL
	if visual, err := s.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	// 2. Try EDITOR
	editor, err := s.envSvc.Editor()
	if err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 3. System-specific primary defaults
	if s.envSvc.IsWindows() {
		return "notepad", nil
	}

	// 4. Common Terminal Editors
	if s.envSvc.IsLinux() || s.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 5. OS Opener Defaults
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
