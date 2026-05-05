package editor

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
)

// EditorResolver defines the interface for resolving the editor command to use.
type EditorResolver interface {
	Resolve(ctx context.Context) (string, error)
}

// DefaultResolver provides the default strategy for resolving an editor.
type DefaultResolver struct {
	envSvc      environment.Service
	cfgProvider ConfigProvider
	explicitCmd string
}

// NewDefaultResolver initializes a new instance of the DefaultResolver.
func NewDefaultResolver(envSvc environment.Service, cfgProvider ConfigProvider, explicitCmd string) *DefaultResolver {
	return &DefaultResolver{
		envSvc:      envSvc,
		cfgProvider: cfgProvider,
		explicitCmd: explicitCmd,
	}
}

// Resolve identifies the editor command based on explicit settings, configuration, environment variables, or OS defaults.
func (r *DefaultResolver) Resolve(ctx context.Context) (string, error) {
	// 1. Explicitly set command
	if strings.TrimSpace(r.explicitCmd) != "" {
		return r.explicitCmd, nil
	}

	// 2. Try configuration
	if r.cfgProvider != nil {
		if cmd, err := r.cfgProvider.GetEditorCommand(ctx); err == nil && strings.TrimSpace(cmd) != "" {
			// TODO: probably a good idea to sanitize this some how
			return cmd, nil
		}
	}

	// 3. Try VISUAL
	if visual, err := r.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		// TODO: maybe a good idea to sanitize all environment vars?
		return visual, nil
	}

	// 4. Try EDITOR
	if editor, err := r.envSvc.Editor(); err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 5. System-specific primary defaults
	if r.envSvc.IsWindows() {
		return "notepad.exe", nil
	}

	// 6. Common Terminal Editors
	if r.envSvc.IsLinux() || r.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 7. OS Opener Defaults
	if r.envSvc.IsMac() {
		if path, err := exec.LookPath("open"); err == nil {
			return path + " -W -t", nil
		}
	}
	if r.envSvc.IsLinux() {
		if path, err := exec.LookPath("xdg-open"); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("could not detect a suitable editor")
}
