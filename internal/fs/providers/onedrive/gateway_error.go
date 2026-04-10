package onedrive

import "fmt"

// GatewayError represents an error that occurred in the GraphDriveGateway.
type GatewayError struct {
	// Operation is the name of the operation that failed.
	Operation string
	// Err is the underlying error.
	Err error
}

// NewGatewayError creates a new GatewayError.
func NewGatewayError(operation string, err error) *GatewayError {
	return &GatewayError{
		Operation: operation,
		Err:       err,
	}
}

// Error returns a formatted error message.
func (e *GatewayError) Error() string {
	return fmt.Sprintf("gateway operation %s failed: %v", e.Operation, e.Err)
}

// Unwrap returns the underlying error.
func (e *GatewayError) Unwrap() error {
	return e.Err
}
