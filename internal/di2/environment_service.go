package di2

type EnvironmentService interface {
	IsWindows() bool
	IsMac() bool
	IsLinux() bool
	ConfigDir() (string, error)
	DataDir() (string, error)
	CacheDir() (string, error)
	LogDir() (string, error)
	InstallDir() (string, error)
	TempDir() (string, error)
	EnsureAll() error
	Name() string
}
