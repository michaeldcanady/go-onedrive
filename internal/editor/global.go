package editor

import (
	"context"
	"fmt"
)

// Context provides the necessary data for state transitions.
type Context struct {
	Session *Session
	Service Service
	Ctx     context.Context
}

var (
	// sessionSM is the global state machine for session lifecycle management.
	sessionSM *StateMachine[State, Event, *Context]
)

func init() {
	sessionSM = NewStateMachine[State, Event, *Context]()

	// Define transitions
	sessionSM.AddTransition([]State{StateCreated}, EventOpen, StateEditing, func(ctx *Context) error {
		// Use type assertion to access implementation-specific methods if necessary
		// Since we're in the same package, we can define a private interface or call unexported methods.
		type editorRunner interface {
			runEditor(ctx context.Context, session *Session) error
		}

		runner, ok := ctx.Service.(editorRunner)
		if !ok {
			return fmt.Errorf("service does not support running an editor")
		}

		if err := runner.runEditor(ctx.Ctx, ctx.Session); err != nil {
			return err
		}

		// Once editor exits, transition to Completed
		return ctx.Session.Handle(ctx.Ctx, ctx.Service, EventComplete)
	})

	sessionSM.AddTransition([]State{StateEditing}, EventComplete, StateCompleted, nil)

	sessionSM.AddTransition([]State{StateCreated, StateEditing, StateCompleted}, EventClose, StateClosed, func(ctx *Context) error {
		type fileRemover interface {
			removeFile(session *Session) error
		}

		remover, ok := ctx.Service.(fileRemover)
		if !ok {
			return fmt.Errorf("service does not support file removal")
		}

		return remover.removeFile(ctx.Session)
	})
}
