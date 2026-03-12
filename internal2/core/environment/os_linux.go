//go:build linux

package environment

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

func localDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".local"), nil
}

func dataBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "share"), nil
}

func cacheBase() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".cache"), nil
}

func stateBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "state"), nil
}

func logsBase() (string, error) {
	state, err := stateBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(state, "logs"), nil
}

func installBase() (string, error) {
	localDir, err := localDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(localDir, "bin"), nil
}

func tempBase() (string, error) {
	return "/tmp", nil
}
