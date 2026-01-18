package environmentservice

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type Service struct {
	appName string
}

func New(appName string) *Service {
	return &Service{appName: appName}
}

func (s *Service) OS() string {
	return runtime.GOOS
}

func (s *Service) IsWindows() bool { return runtime.GOOS == "windows" }
func (s *Service) IsMac() bool     { return runtime.GOOS == "darwin" }
func (s *Service) IsLinux() bool   { return runtime.GOOS == "linux" }

func (s *Service) ConfigDir(ctx context.Context) (string, error) {
	base, err := configBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *Service) DataDir(ctx context.Context) (string, error) {
	base, err := dataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *Service) CacheDir(ctx context.Context) (string, error) {
	base, err := cacheBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, s.appName), nil
}

func (s *Service) LogDir(ctx context.Context) (string, error) {
	base, err := logsBase()
	if err != nil {
		return "", err
	}

	if s.IsLinux() {
		// ~/.local/state/<app>/logs
		return filepath.Join(base, s.appName, "logs"), nil
	}

	// Windows: %LOCALAPPDATA%\Logs\<app>
	// macOS:   ~/Library/Logs/<app>
	return filepath.Join(base, s.appName), nil
}

func (s *Service) InstallDir(ctx context.Context) (string, error) {
	base, err := installBase()
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

func (s *Service) TempDir(ctx context.Context) (string, error) {
	return tempBase()
}

func (s *Service) EnsureAll(ctx context.Context) error {
	// Don’t create InstallDir or TempDir
	creators := []func(context.Context) (string, error){
		s.ConfigDir,
		s.DataDir,
		s.CacheDir,
		s.LogDir,
	}

	for _, fn := range creators {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		dir, err := fn(ctx)
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
