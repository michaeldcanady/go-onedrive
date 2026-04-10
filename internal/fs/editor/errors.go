package editor

import "fmt"

// EditorLaunchError indicates that the external editor could not be started or failed during execution.
type EditorLaunchError struct {
	// Path is the path of the file that was being edited.
	Path string
	// Err is the underlying error.
	Err error
}

func NewEditorLaunchError(path string, err error) *EditorLaunchError {
	return &EditorLaunchError{
		Path: path,
		Err:  err,
	}
}

func (e *EditorLaunchError) Error() string {
	return fmt.Sprintf("failed to launch editor for path %s: %v", e.Path, e.Err)
}

func (e *EditorLaunchError) Unwrap() error {
	return e.Err
}

// NoEditorDetectedError indicates that no suitable editor was found in the environment.
type NoEditorDetectedError struct {
	// Err is the underlying error, if any.
	Err error
}

func NewNoEditorDetectedError(err error) *NoEditorDetectedError {
	return &NoEditorDetectedError{
		Err: err,
	}
}

func (e *NoEditorDetectedError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("could not detect a suitable editor: %v", e.Err)
	}
	return "could not detect a suitable editor"
}

func (e *NoEditorDetectedError) Unwrap() error {
	return e.Err
}
