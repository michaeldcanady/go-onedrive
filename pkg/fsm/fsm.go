package fsm

import (
	"context"
)

// State defines a single step in a state machine.
// It takes a context and a pointer to the machine's shared data.
// It returns the next state to transition to, or nil if the machine should stop.
type State[T any] interface {
	Execute(ctx context.Context, data *T) (State[T], error)
}

// StateFunc is an adapter to allow the use of ordinary functions as states.
type StateFunc[T any] func(ctx context.Context, data *T) (State[T], error)

// Execute calls f(ctx, data).
func (f StateFunc[T]) Execute(ctx context.Context, data *T) (State[T], error) {
	return f(ctx, data)
}

// Machine runs a sequence of states until a nil state is returned or an error occurs.
type Machine[T any] struct {
	data *T
}

// NewMachine initializes a new state machine with the provided shared data.
func NewMachine[T any](data *T) *Machine[T] {
	return &Machine[T]{data: data}
}

// Run starts the state machine with the given initial state.
func (m *Machine[T]) Run(ctx context.Context, initial State[T]) error {
	current := initial
	for current != nil {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			next, err := current.Execute(ctx, m.data)
			if err != nil {
				return err
			}
			current = next
		}
	}
	return nil
}
