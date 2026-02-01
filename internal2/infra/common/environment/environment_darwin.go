//go:build darwin

package environment

import (
	"os"
	"path/filepath"
)

func ConfigBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Preferences"), nil
}

func libraryBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Library"), nil
}

func applicationSupportBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Application Support"), nil
}

func StateBase() (string, error) {
	return applicationSupportBase()
}

func DataBase() (string, error) {
	return applicationSupportBase()
}

func CacheBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Caches"), nil
}

func LogsBase() (string, error) {
	libraryDir, err := libraryBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(libraryDir, "Logs"), nil
}

func InstallBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, "Applications"), nil
}

func TempBase() (string, error) {
	return os.TempDir(), nil
}
