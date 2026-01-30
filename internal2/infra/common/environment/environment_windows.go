//go:build windows

package environment

import (
	"fmt"
	"os"
	"path/filepath"
)

func ConfigBase() (string, error) {
	base := os.Getenv("APPDATA")
	if base == "" {
		return "", fmt.Errorf("APPDATA not set")
	}
	return base, nil
}

func DataBase() (string, error) {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" {
		return "", fmt.Errorf("LOCALAPPDATA not set")
	}
	return base, nil
}

func CacheBase() (string, error) {
	base, err := DataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Cache"), nil
}

func LogsBase() (string, error) {
	base, err := DataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Logs"), nil
}

func InstallBase() (string, error) {
	base, err := DataBase()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "Programs"), nil
}

func TempBase() (string, error) {
	return os.TempDir(), nil
}
