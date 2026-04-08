//go:build darwin

package environment

import (
	"os"
	"path/filepath"
)

// configBase returns the base directory for configuration files on macOS.
func configBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Preferences"), nil
}

// libraryBase returns the path to the user's Library directory on macOS.
func libraryBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library"), nil
}

// applicationSupportBase returns the path to the Application Support directory on macOS.
func applicationSupportBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Application Support"), nil
}

func stateBase() (string, error) {
	return applicationSupportBase()
}

func dataBase() (string, error) {
	return applicationSupportBase()
}

func cacheBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Caches"), nil
}

func logsBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Logs"), nil
}

func installBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Applications"), nil
}

func tempBase() (string, error) {
	return os.TempDir(), nil
}
