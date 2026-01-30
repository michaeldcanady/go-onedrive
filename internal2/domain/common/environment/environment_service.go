package environment

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
}
