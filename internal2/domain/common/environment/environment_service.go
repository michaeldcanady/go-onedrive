package environment

import "github.com/michaeldcanady/go-onedrive/internal2/infra/common/logging"

type EnvironmentService interface {
	CacheDir() (string, error)
	ConfigDir() (string, error)
	DataDir() (string, error)
	EnsureAll() error
	InstallDir() (string, error)
	IsLinux() bool
	IsMac() bool
	IsWindows() bool
	LogDir() (string, error)
	Name() string
	OS() string
	TempDir() (string, error)
	StateDir() (string, error)
	OutputDestination() (logging.OutputDestination, error)
	LogLevel() (string, error)
	SHELL() (string, error)
}
