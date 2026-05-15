package editor

import (
	"context"
	"os"
	"os/exec"
)

// Service coordinates the launching of external editor processes.
type Service interface {
	// Open launches the system's preferred editor (from the EDITOR environment variable)
	// to modify the file at the specified path. It blocks until the editor process exits.
	Open(ctx context.Context, path string) error
}

type editorService struct{}

// NewEditorService returns a new [Service] implementation.
func NewEditorService() Service {
	return &editorService{}
}

func (s *editorService) Open(ctx context.Context, path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // Fallback
	}

	executable, err := exec.LookPath(editor)
	if err != nil {
		return err
	}

	// nolint:gosec
	cmd := exec.CommandContext(ctx, executable, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
