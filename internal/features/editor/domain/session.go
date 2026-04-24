package editor

import (
	"context"
	"fmt"

	fs "github.com/michaeldcanady/go-onedrive/internal/features/fs/domain"
)

// Session represents the state of an editing lifecycle for a single file.
type Session struct {
	// ID is the unique identifier for this session.
	ID string
	// RemoteURI is the original location of the file being edited.
	RemoteURI *fs.URI
	// LocalURI is the temporary local path where the file is staged.
	LocalURI *fs.URI
	// InitialHash is the SHA256 hash of the content when the session was created.
	InitialHash []byte

	state State
}

// State returns the current lifecycle phase of the session.
func (s *Session) State() State {
	return s.state
}

// SetState updates the current lifecycle phase of the session.
func (s *Session) SetState(state State) {
	s.state = state
}

// Handle triggers a state transition on the session based on the provided event.
func (s *Session) Handle(ctx context.Context, svc Service, event Event) error {
	ds, ok := svc.(*DefaultService)
	if !ok {
		return fmt.Errorf("unsupported service type")
	}

	return ds.sm.Handle(s, event, &Context{
		Session: s,
		Service: svc,
		Ctx:     ctx,
	})
}
