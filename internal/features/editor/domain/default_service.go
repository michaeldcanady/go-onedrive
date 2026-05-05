package editor

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"

	"github.com/google/shlex"
	environment "github.com/michaeldcanady/go-onedrive/internal/core/env"
	"github.com/michaeldcanady/go-onedrive/internal/core/logger"
	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
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

// WithResolver sets the editor resolver for the service.
func WithResolver(resolver EditorResolver) Option {
	return func(s *DefaultService) {
		s.resolver = resolver
	}
}

// WithSessionManager sets the session manager for the service.
func WithSessionManager(sm SessionManager) Option {
	return func(s *DefaultService) {
		s.sessions = sm
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
	resolver    EditorResolver
	sessions    SessionManager
	sm          *StateMachine[State, Event, *Context]
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

	if s.resolver == nil {
		s.resolver = NewDefaultResolver(s.envSvc, s.cfgProvider, s.editorCmd)
	}

	if s.sessions == nil {
		s.sessions = NewDefaultSessionManager(s.envSvc, s.uriFactory)
	}

	s.sm = s.setupStateMachine()

	return s
}

func (s *DefaultService) setupStateMachine() *StateMachine[State, Event, *Context] {
	sm := NewStateMachine[State, Event, *Context]()

	// Define transitions
	sm.AddTransition([]State{StateCreated}, EventOpen, StateEditing, nil)
	sm.AddTransition([]State{StateEditing}, EventComplete, StateCompleted, nil)

	sm.AddTransition([]State{StateCreated, StateEditing, StateCompleted}, EventClose, StateClosed, func(ctx *Context) error {
		type fileRemover interface {
			removeFile(session *Session) error
		}

		remover, ok := ctx.Service.(fileRemover)
		if !ok {
			return fmt.Errorf("service does not support file removal")
		}
		return remover.removeFile(ctx.Session)
	})

	return sm
}

func (s *DefaultService) removeFile(session *Session) error {
	type internalRemover interface {
		removeFile(session *Session) error
	}
	if ir, ok := s.sessions.(internalRemover); ok {
		return ir.removeFile(session)
	}
	return fmt.Errorf("session manager does not support internal file removal")
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
		resolver:    s.resolver,
		sessions:    s.sessions,
	}

	for _, opt := range opts {
		opt(newS)
	}

	// Re-initialize resolver if editorCmd changed and it's a DefaultResolver
	if _, ok := newS.resolver.(*DefaultResolver); ok {
		newS.resolver = NewDefaultResolver(newS.envSvc, newS.cfgProvider, newS.editorCmd)
	}

	newS.sm = newS.setupStateMachine()

	return newS
}

// CreateSession initializes a new editing session.
func (s *DefaultService) CreateSession(ctx context.Context, remoteURI *fs.URI, r io.Reader) (*Session, error) {
	return s.sessions.CreateSession(ctx, remoteURI, r)
}

// Modified checks if the local file in the session has changed.
func (s *DefaultService) Modified(session *Session) (bool, error) {
	return s.sessions.Modified(session)
}

// NewContent returns a reader for the modified content in the session.
func (s *DefaultService) NewContent(session *Session) (io.ReadCloser, error) {
	return s.sessions.NewContent(session)
}

// Cleanup removes the temporary local file and releases session resources.
func (s *DefaultService) Cleanup(ctx context.Context, session *Session) error {
	return s.sessions.Cleanup(ctx, s, session)
}

func (s *DefaultService) getEditorParts(ctx context.Context) ([]string, error) {
	editorCmd, err := s.resolver.Resolve(ctx)
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
	if err := session.Handle(ctx, s, EventOpen); err != nil {
		return err
	}

	if err := s.runEditor(ctx, session); err != nil {
		return err
	}

	return session.Handle(ctx, s, EventComplete)
}

// runEditor is the internal implementation that actually launches the editor.
func (s *DefaultService) runEditor(ctx context.Context, session *Session) error {
	editorParts, err := s.getEditorParts(ctx)
	if err != nil {
		return err
	}

	args := append(editorParts, session.LocalURI.Path)
	// nolint:gosec // 204 // doesn't allow for user input
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
