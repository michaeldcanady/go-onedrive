package edit

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

const (
	defaultUnixShell     = "/bin/bash"
	defaultUnixEditor    = "vim"
	defaultWindowsShell  = "cmd"
	defaultWindowsEditor = "notepad"
)

type EditorService struct {
	logger logging.Logger
	envSvc environment.EnvironmentService

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	editor string
}

func NewEditorService(env environment.EnvironmentService, logger logging.Logger) *EditorService {
	return &EditorService{
		envSvc: env,
		logger: logger,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

func (s *EditorService) WithIO(stdin io.Reader, stdout, stderr io.Writer) *EditorService {
	s.Stdin = stdin
	s.Stdout = stdout
	s.Stderr = stderr
	return s
}

func (s *EditorService) getEditor() (string, error) {
	if s.editor != "" {
		return s.editor, nil
	}

	if s.envSvc.IsWindows() {
		return defaultWindowsEditor, nil
	}

	editor, err := s.envSvc.Editor()
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(editor) == "" {
		editor = defaultUnixEditor
	}

	return editor, nil
}

func (s *EditorService) getShell() ([]string, error) {
	if s.envSvc.IsWindows() {
		return []string{defaultWindowsShell, "/C"}, nil
	}

	shell, err := s.envSvc.Shell()
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(shell) == "" {
		shell = defaultUnixShell
	}

	return []string{shell, "-c"}, nil
}

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
