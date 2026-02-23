package edit

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

const (
	defaultUnixShell    = "/bin/bash"
	defaultWindowsShell = "cmd"
)

// Editor defines the interface for launching an external editor.
type Editor interface {
	Launch(path string) error
	LaunchTempFile(prefix, suffix string, reader io.Reader) ([]byte, string, error)
}

// EditorService provides functionality to launch an external editor.
type EditorService struct {
	logger logging.Logger
	envSvc environment.EnvironmentService

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	editor string
}

// NewEditorService creates a new instance of EditorService.
func NewEditorService(env environment.EnvironmentService, logger logging.Logger) *EditorService {
	return &EditorService{
		envSvc: env,
		logger: logger,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// WithIO sets the input and output streams for the editor process.
func (s *EditorService) WithIO(stdin io.Reader, stdout, stderr io.Writer) *EditorService {
	s.Stdin = stdin
	s.Stdout = stdout
	s.Stderr = stderr
	return s
}

// getEditor determines the editor command to use based on environment variables and OS defaults.
func (s *EditorService) getEditor() (string, error) {
	if s.editor != "" {
		return s.editor, nil
	}

	// 1. Try VISUAL (standard Unix preference for full-screen editors)
	if visual, err := s.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	// 2. Try EDITOR (standard Unix preference for editors)
	editor, err := s.envSvc.Editor()
	if err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 3. System-specific primary defaults
	if s.envSvc.IsWindows() {
		return "notepad", nil
	}

	// 4. Common Terminal Editors (preferred for CLI tools on Unix-like systems)
	if s.envSvc.IsLinux() || s.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 5. OS Opener Defaults (last resort, might try to launch a GUI)
	if s.envSvc.IsMac() {
		if path, err := exec.LookPath("open"); err == nil {
			return path + " -W -t", nil // -W wait, -t use text editor
		}
	}
	if s.envSvc.IsLinux() {
		if path, err := exec.LookPath("xdg-open"); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not detect a suitable editor")
}

// getShell determines the shell command to use for launching the editor.
func (s *EditorService) getShell() ([]string, error) {
	if s.envSvc.IsWindows() {
		return []string{defaultWindowsShell, "/C"}, nil
	}

	shell, err := s.envSvc.Shell()
	if err != nil || strings.TrimSpace(shell) == "" {
		shell = defaultUnixShell
	}

	return []string{shell, "-c"}, nil
}

// LaunchTempFile creates a temporary file with the given content, launches the editor,
// and returns the updated content after the editor closes.
func (s *EditorService) LaunchTempFile(prefix, suffix string, reader io.Reader) ([]byte, string, error) {
	if !strings.HasPrefix(suffix, ".") {
		suffix = fmt.Sprintf(".%s", suffix)
	}
	f, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	path := f.Name()
	if _, err := io.Copy(f, reader); err != nil {
		os.Remove(path)
		return nil, path, err
	}

	f.Close()

	if err := s.Launch(path); err != nil {
		return nil, path, err
	}

	bytes, err := os.ReadFile(path)
	return bytes, path, err
}

// buildCmd constructs the command arguments for launching the editor.
func (s *EditorService) buildCmd(path string) ([]string, error) {
	shell, err := s.getShell()
	if err != nil {
		return nil, err
	}

	editor, err := s.getEditor()
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s %s", editor, path)
	return append(shell, cmd), nil
}

// Launch executes the editor command and waits for it to exit.
func (s *EditorService) Launch(path string) error {
	args, err := s.buildCmd(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr

	return cmd.Run()
}
