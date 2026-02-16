//go:build !windows

package edit

import (
	"fmt"
	"os/exec"
)

func (s *EditorService) buildCmd(path string) ([]string, error) {
	shell, err := s.getShell()
	if err != nil {
		return nil, err
	}

	editor, err := s.getEditor()
	if err != nil {
		return nil, err
	}

	cmd := fmt.Sprintf("%s %s", editor, path)
	return append(shell, cmd), nil
}

func (s *EditorService) Launch(path string) error {
	args, err := s.buildCmd(path)
	if err != nil {
		return err
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = s.Stdin
	cmd.Stdout = s.Stdout
	cmd.Stderr = s.Stderr

	return cmd.Run()
}
