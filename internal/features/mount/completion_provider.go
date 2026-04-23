package mount

import "context"

// CompletionProvider defines the interface for backends that support dynamic completion.
type CompletionProvider interface {
	GetOptionCompletions(ctx context.Context, identityID string, toComplete string) ([]string, error)
}
