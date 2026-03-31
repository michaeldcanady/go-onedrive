package environment

// Service defines the interface for environment-related operations.
type Service interface {
	// CacheDir returns the path to the cache directory.
	CacheDir() (string, error)
	// ConfigDir returns the path to the configuration directory.
	ConfigDir() (string, error)
	// DataDir returns the path to the data directory.
	DataDir() (string, error)
	// EnsureAll ensures all necessary directories exist.
	EnsureAll() error
	// InstallDir returns the path to the installation directory.
	InstallDir() (string, error)
	// IsLinux returns true if the OS is Linux.
	IsLinux() bool
	// IsMac returns true if the OS is macOS.
	IsMac() bool
	// IsWindows returns true if the OS is Windows.
	IsWindows() bool
	// LogDir returns the path to the log directory.
	LogDir() (string, error)
	// Name returns the name of the application.
	Name() string
	// OS returns the name of the operating system.
	OS() string
	// TempDir returns the path to the temporary directory.
	TempDir() (string, error)
	// StateDir returns the path to the state directory.
	StateDir() (string, error)
	// Shell returns the preferred shell.
	Shell() (string, error)
	// Editor returns the preferred editor.
	Editor() (string, error)
	// Visual returns the preferred visual editor.
	Visual() (string, error)
}
