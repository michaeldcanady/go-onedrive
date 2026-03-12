//go:build windows

package environment

import (
	"fmt"
	"os"
	"path/filepath"
)

func configBase() (string, error) {
	base := os.Getenv("APPDATA")
	if base == "" {
		return "", fmt.Errorf("APPDATA not set")
	}
	return base, nil
}

func localAppdata() (string, error) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return "", fmt.Errorf("LOCALAPPDATA not set")
	}
	return base, nil
}

func stateBase() (string, error) {
	base, err := localAppdata()
	if err != nil {
		return "", err
	}
	return base, nil
}

func dataBase() (string, error) {
	base, err := localAppdata()
	if err != nil {
		return "", err
	}
	return base, nil
}

func cacheBase() (string, error) {
	base, err := dataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Cache"), nil
}

func logsBase() (string, error) {
	base, err := dataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Logs"), nil
}

func installBase() (string, error) {
	base, err := dataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Programs"), nil
}

func tempBase() (string, error) {
	return os.TempDir(), nil
}
