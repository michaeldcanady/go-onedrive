package editor

import (
	"fmt"
	"io"
	"os"
	"strings"

	domaineditor "github.com/michaeldcanady/go-onedrive/internal2/domain/editor"
	"github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"
)

var _ domaineditor.Service = (*Service)(nil)

// Service coordinates temporary file management and editor launching.
type Service struct {
	launcher domaineditor.Launcher
	logger   logging.Logger
}

func NewService(launcher domaineditor.Launcher, logger logging.Logger) *Service {
	return &Service{
		launcher: launcher,
		logger:   logger,
	}
}

func (s *Service) WithIO(stdin io.Reader, stdout, stderr io.Writer) domaineditor.Service {
	s.launcher.WithIO(stdin, stdout, stderr)
	return s
}

// LaunchTempFile creates a temporary file with the given content, launches the editor,
// and returns the updated content after the editor closes.
func (s *Service) LaunchTempFile(prefix, suffix string, reader io.Reader) ([]byte, string, error) {
	if !strings.HasPrefix(suffix, ".") {
		suffix = fmt.Sprintf(".%s", suffix)
	}

	f, err := os.CreateTemp("", prefix+"*"+suffix)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	path := f.Name()
	if _, err := io.Copy(f, reader); err != nil {
		os.Remove(path)
		return nil, path, err
	}

	f.Close()

	if err := s.launcher.Launch(path); err != nil {
		return nil, path, err
	}

	bytes, err := os.ReadFile(path)
	return bytes, path, err
}
