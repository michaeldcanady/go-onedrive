//go:build linux

package environment

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

func localDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local"), nil
}

func DataBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "share"), nil
}

func CacheBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

func StateBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "state"), nil
}

func LogsBase() (string, error) {
	state, err := StateBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(state, "logs"), nil
}

func InstallBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "bin"), nil
}

func TempBase() (string, error) {
	return "/tmp", nil
}
