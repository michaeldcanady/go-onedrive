package editor

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	domainenvironment "github.com/michaeldcanady/go-onedrive/internal2/domain/common/environment"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/common/logger"
	"github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
)

const (
	defaultUnixShell    = "/bin/bash"
	defaultWindowsShell = "cmd"
)

var _ editor.Launcher = (*Launcher)(nil)

// Launcher implements domain.editor.Launcher using os/exec.
type Launcher struct {
	envSvc domainenvironment.EnvironmentService
	log    logger.Logger
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
}

func NewLauncher(envSvc domainenvironment.EnvironmentService, l logger.Logger) *Launcher {
	return &Launcher{
		envSvc: envSvc,
		log:    l,
		stdin:  os.Stdin,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

func (l *Launcher) WithIO(stdin io.Reader, stdout, stderr io.Writer) editor.Launcher {
	l.stdin = stdin
	l.stdout = stdout
	l.stderr = stderr
	return l
}

func (l *Launcher) Launch(path string) error {
	args, err := l.buildCmd(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = l.stdin
	cmd.Stdout = l.stdout
	cmd.Stderr = l.stderr

	return cmd.Run()
}

func (l *Launcher) buildCmd(path string) ([]string, error) {
	shell, err := l.getShell()
	if err != nil {
		return nil, err
	}

	editorCmd, err := l.getEditor()
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s %s", editorCmd, path)
	return append(shell, cmd), nil
}

func (l *Launcher) getEditor() (string, error) {
	// 1. Try VISUAL
	if visual, err := l.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	// 2. Try EDITOR
	editor, err := l.envSvc.Editor()
	if err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 3. System-specific primary defaults
	if l.envSvc.IsWindows() {
		return "notepad", nil
	}

	// 4. Common Terminal Editors
	if l.envSvc.IsLinux() || l.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 5. OS Opener Defaults
	if l.envSvc.IsMac() {
		if path, err := exec.LookPath("open"); err == nil {
			return path + " -W -t", nil
		}
	}
	if l.envSvc.IsLinux() {
		if path, err := exec.LookPath("xdg-open"); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not detect a suitable editor")
}

func (l *Launcher) getShell() ([]string, error) {
	if l.envSvc.IsWindows() {
		return []string{defaultWindowsShell, "/C"}, nil
	}

	shell, err := l.envSvc.Shell()
	if err != nil || strings.TrimSpace(shell) == "" {
		shell = defaultUnixShell
	}

	return []string{shell, "-c"}, nil
}
