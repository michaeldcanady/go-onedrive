//go:build linux

package environmentservice

import (
	"os"
	"path/filepath"
)

func configBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config"), nil
}

func dataBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "share"), nil
}

func cacheBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

func logsBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// ~/.local/state/yourapp/logs (appName is added in impl.go)
	return filepath.Join(home, ".local", "state"), nil
}

func installBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local", "bin"), nil
}

func tempBase() (string, error) {
	return "/tmp", nil
}
