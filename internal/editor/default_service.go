package editor

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
	"github.com/google/uuid"
	fs "github.com/michaeldcanady/go-onedrive/internal/core/fs"
	"github.com/michaeldcanady/go-onedrive/internal/features/environment"
	"github.com/michaeldcanady/go-onedrive/internal/features/logger"
)

// ConfigProvider defines the interface required to fetch editor configuration.
type ConfigProvider interface {
	GetEditorCommand(ctx context.Context) (string, error)
}

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

// WithConfig sets the configuration provider for the editor.
func WithConfig(cfgProvider ConfigProvider) Option {
	return func(s *DefaultService) {
		s.cfgProvider = cfgProvider
	}
}

// DefaultService provides the default implementation of the editor service.
type DefaultService struct {
	envSvc      environment.Service
	cfgProvider ConfigProvider
	uriFactory  *fs.URIFactory
	log         logger.Logger
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
	editorCmd   string
}

// NewDefaultService initializes a new instance of the DefaultService.
func NewDefaultService(envSvc environment.Service, uriFactory *fs.URIFactory, l logger.Logger, opts ...Option) *DefaultService {
	s := &DefaultService{
		envSvc:     envSvc,
		uriFactory: uriFactory,
		log:        l,
		stdin:      os.Stdin,
		stdout:     os.Stdout,
		stderr:     os.Stderr,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

// WithOptions returns a new Service instance with the specified options applied.
func (s *DefaultService) WithOptions(opts ...Option) Service {
	newS := &DefaultService{
		envSvc:      s.envSvc,
		cfgProvider: s.cfgProvider,
		uriFactory:  s.uriFactory,
		log:         s.log,
		stdin:       s.stdin,
		stdout:      s.stdout,
		stderr:      s.stderr,
		editorCmd:   s.editorCmd,
	}

	for _, opt := range opts {
		opt(newS)
	}

	return newS
}

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

func (s *DefaultService) getEditorCmd() (string, error) {
	// 1. Explicitly set command
	if strings.TrimSpace(s.editorCmd) != "" {
		return s.editorCmd, nil
	}

	// 2. Try configuration
	if s.cfgProvider != nil {
		if cmd, err := s.cfgProvider.GetEditorCommand(context.Background()); err == nil && strings.TrimSpace(cmd) != "" {
			return cmd, nil
		}
	}

	// 3. Try VISUAL
	if visual, err := s.envSvc.Visual(); err == nil && strings.TrimSpace(visual) != "" {
		return visual, nil
	}

	// 4. Try EDITOR
	if editor, err := s.envSvc.Editor(); err == nil && strings.TrimSpace(editor) != "" {
		return editor, nil
	}

	// 5. System-specific primary defaults
	if s.envSvc.IsWindows() {
		return "notepad.exe", nil
	}

	// 6. Common Terminal Editors
	if s.envSvc.IsLinux() || s.envSvc.IsMac() {
		fallbacks := []string{"vim", "vi", "nano"}
		for _, f := range fallbacks {
			if path, err := exec.LookPath(f); err == nil {
				return path, nil
			}
		}
	}

	// 7. OS Opener Defaults
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

// Open launches the configured editor for the given session and waits for it to exit.
func (s *DefaultService) Open(ctx context.Context, session *Session) error {
	return session.Handle(ctx, s, EventOpen)
}

// runEditor is the internal implementation that actually launches the editor.
func (s *DefaultService) runEditor(ctx context.Context, session *Session) error {
	editorParts, err := s.getEditorParts()
	if err != nil {
		return err
	}

	args := append(editorParts, session.LocalURI.Path)
	cmd := exec.CommandContext(ctx, args[0], args[1:]...)
	cmd.Stdin = s.stdin
	cmd.Stdout = s.stdout
	cmd.Stderr = s.stderr

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start)

	if err != nil {
		return fmt.Errorf("editor failed: %w", err)
	}

	// If the editor returned very quickly (less than 1 second), it likely backgrounded.
	if duration < 1*time.Second {
		s.log.Warn("editor returned very quickly; if it opened in the background, ensure it was saved before continuing",
			logger.Duration("duration", duration))
	}

	return nil
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
