package fs

import "fmt"

// Validator defines the interface for options that can be validated.
type Validator interface {
	Validate() error
}

// CopyOptions defines the configuration for a copy operation.
type CopyOptions struct {
	// Recursive determines whether to include nested items.
	Recursive bool
	// Overwrite determines whether to replace existing items at the destination.
	Overwrite bool
}

func (o CopyOptions) Validate() error {
	return nil
}

// ListOptions defines the configuration for an enumeration operation.
type ListOptions struct {
	// Recursive determines whether to traverse into subdirectories.
	Recursive bool
}

func (o ListOptions) Validate() error {
	return nil
}

// ReadOptions defines the configuration for a file reading operation.
type ReadOptions struct{}

func (o ReadOptions) Validate() error {
	return nil
}

// WriteOptions defines the configuration for a write operation.
type WriteOptions struct {
	// Overwrite determines whether to replace an existing item.
	Overwrite bool
	// IfMatch is the ETag of the item that should be overwritten.
	IfMatch string
	// Size is the total number of bytes to be written (required for resumable uploads).
	Size int64
}

func (o WriteOptions) Validate() error {
	if o.Size < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	return nil
}
