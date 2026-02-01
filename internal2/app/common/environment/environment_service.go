package environment

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	infracommon "github.com/michaeldcanady/go-onedrive/internal2/infra/common/environment"
)

type EnvironmentService struct {
	appName string
}

func New2(appName string) *EnvironmentService {
	return &EnvironmentService{appName: appName}
}

func (s *EnvironmentService) OS() string {
	return runtime.GOOS
}

func (s *EnvironmentService) Name() string {
	return s.appName
}

func (s *EnvironmentService) IsWindows() bool { return runtime.GOOS == "windows" }
func (s *EnvironmentService) IsMac() bool     { return runtime.GOOS == "darwin" }
func (s *EnvironmentService) IsLinux() bool   { return runtime.GOOS == "linux" }

func (s *EnvironmentService) ConfigDir() (string, error) {
	base, err := infracommon.ConfigBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) DataDir() (string, error) {
	base, err := infracommon.DataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) CacheDir() (string, error) {
	base, err := infracommon.CacheBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) LogDir() (string, error) {
	if s.IsLinux() {
		base, err := infracommon.StateBase()
		if err != nil {
			return "", err
		}

		return filepath.Join(base, s.appName, "logs"), nil
	}
	base, err := infracommon.LogsBase()
	if err != nil {
		return "", err
	}

	// Windows: %LOCALAPPDATA%\Logs\<app>
	// macOS:   ~/Library/Logs/<app>
	return filepath.Join(base, s.appName), nil
}

func (s *EnvironmentService) InstallDir() (string, error) {
	base, err := infracommon.InstallBase()
	if err != nil {
		return "", err
	}

	if s.IsWindows() {
		// %LOCALAPPDATA%\Programs\<app>
		return filepath.Join(base, s.appName), nil
	}

	if s.IsLinux() {
		// ~/.local/bin (binary itself lives here, not a subdirectory)
		return base, nil
	}

	// macOS: ~/Applications (you may later decide on ~/Applications/<app>.app)
	return base, nil
}

func (s *EnvironmentService) TempDir() (string, error) {
	temp, err := infracommon.TempBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(temp, s.appName), nil
}

func (s *EnvironmentService) StateDir() (string, error) {
	state, err := infracommon.StateBase()
	if err != nil {
		return "", err
	}

	return filepath.Join(state, s.appName), nil
}

func (s *EnvironmentService) EnsureAll() error {
	// Don’t create InstallDir or TempDir
	creators := []func() (string, error){
		s.ConfigDir,
		s.DataDir,
		s.CacheDir,
		s.LogDir,
		s.StateDir,
	}

	for _, fn := range creators {
		dir, err := fn()
		if err != nil {
			return err
		}
		if dir == "" {
			return fmt.Errorf("empty directory path resolved")
		}
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}

	return nil
}
