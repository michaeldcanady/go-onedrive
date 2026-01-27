//go:build linux

package common

import (
	"os"
	"path/filepath"
)

func ConfigBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

func DataBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}

func CacheBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

func LogsBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// ~/.local/state/yourapp/logs (appName is added in impl.go)
	return filepath.Join(home, ".local", "state"), nil
}

func InstallBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin"), nil
}

func TempBase() (string, error) {
	return "/tmp", nil
}
