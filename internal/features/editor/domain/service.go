package editor

import (
	"context"
	"fmt"
	"io"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// State represents the current phase of a session's lifecycle.
type State int

const (
	// StateCreated indicates the session has been initialized and the local file is staged.
	StateCreated State = iota
	// StateEditing indicates the external editor is currently open.
	StateEditing
	// StateCompleted indicates the editor has exited.
	StateCompleted
	// StateClosed indicates the session has been cleaned up and resources released.
	StateClosed
)

// String returns the string representation of the state.
func (s State) String() string {
	switch s {
	case StateCreated:
		return "Created"
	case StateEditing:
		return "Editing"
	case StateCompleted:
		return "Completed"
	case StateClosed:
		return "Closed"
	default:
		return fmt.Sprintf("Unknown(%d)", s)
	}
}

// Event represents an action that triggers a state transition.
type Event int

const (
	// EventOpen is triggered when the editor is opened.
	EventOpen Event = iota
	// EventComplete is triggered when the editor exits.
	EventComplete
	// EventClose is triggered when the session is cleaned up.
	EventClose
)

// String returns the string representation of the event.
func (e Event) String() string {
	switch e {
	case EventOpen:
		return "Open"
	case EventComplete:
		return "Complete"
	case EventClose:
		return "Close"
	default:
		return fmt.Sprintf("Unknown(%d)", e)
	}
}

// Context provides the necessary data for state transitions.
type Context struct {
	Session *Session
	Service Service
	Ctx     context.Context
}

// SessionManager defines the interface for managing the lifecycle of editing sessions.
type SessionManager interface {
	CreateSession(ctx context.Context, remoteURI *fs.URI, reader io.Reader) (*Session, error)
	Modified(session *Session) (bool, error)
	NewContent(session *Session) (io.ReadCloser, error)
	Cleanup(ctx context.Context, svc Service, session *Session) error
}

// Service defines the interface for editor-related operations and session management.
type Service interface {
	// CreateSession initializes a new editing session for the given remote URI.
	CreateSession(ctx context.Context, remoteURI *fs.URI, reader io.Reader) (*Session, error)

	// WithOptions returns a new Service instance with the specified options applied.
	WithOptions(opts ...Option) Service

	// Open launches the configured editor for the given session and waits for it to exit.
	Open(ctx context.Context, session *Session) error

	// Modified checks if the local file in the session has changed since it was staged.
	Modified(session *Session) (bool, error)

	// NewContent returns a reader for the modified content in the session.
	NewContent(session *Session) (io.ReadCloser, error)

	// Cleanup removes the temporary local file and releases session resources.
	Cleanup(ctx context.Context, session *Session) error
}
