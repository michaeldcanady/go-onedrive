package editor

import (
	"fmt"
	"sync"
)

// Stateful is an interface for objects that have a state.
type Stateful[S comparable] interface {
	State() S
	SetState(S)
}

// Action represents a function to be executed during a transition.
type Action[T any] func(T) error

type tuple[T1, T2 any] struct {
	first  T1
	second T2
}

// StateMachine manages the lifecycle transitions of a generic entity.
type StateMachine[S comparable, E comparable, T any] struct {
	mu          sync.RWMutex
	transitions map[tuple[S, E]]tuple[S, Action[T]]
}

// NewStateMachine initializes a new StateMachine.
func NewStateMachine[S comparable, E comparable, T any]() *StateMachine[S, E, T] {
	return &StateMachine[S, E, T]{
		transitions: make(map[tuple[S, E]]tuple[S, Action[T]]),
	}
}

// AddTransition registers a valid transition from one or more source states to a target state triggered by an event.
func (s *StateMachine[S, E, T]) AddTransition(fromStates []S, event E, toState S, action Action[T]) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, from := range fromStates {
		s.transitions[tuple[S, E]{first: from, second: event}] = tuple[S, Action[T]]{first: toState, second: action}
	}
}

// Handle triggers a state transition on the stateful object based on the provided event and context.
func (s *StateMachine[S, E, T]) Handle(stateful Stateful[S], event E, context T) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fromState := stateful.State()
	key := tuple[S, E]{first: fromState, second: event}
	next, ok := s.transitions[key]
	if !ok {
		return fmt.Errorf("invalid transition: state %v + event %v", fromState, event)
	}

	toState, action := next.first, next.second

	// Set state before action to allow actions to see the new state (e.g., "Editing")
	stateful.SetState(toState)

	if action != nil {
		if err := action(context); err != nil {
			// Rollback state on failure?
			// For now, we just return the error.
			return fmt.Errorf("action failed during transition: %w", err)
		}
	}

	return nil
}
