// Package paths provides convenience functions related to the file system.
package paths

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/adrg/xdg"

	"github.com/eikendev/hackenv/internal/constants"
)

// GetDataFilePath returns a file from the XDG data directory.
func GetDataFilePath(file string) (string, error) {
	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, file))
	if err != nil {
		slog.Error("Cannot access data directory", "file", file, "err", err)
		return "", fmt.Errorf("cannot resolve data file %q: %w", file, err)
	}

	return path, nil
}

// EnsureDirExists creates the given directory if it does not exists.
func EnsureDirExists(path string) error {
	err := os.MkdirAll(path, 0o750)
	if err != nil {
		slog.Error("Cannot create directory", "path", path, "err", err)
		return fmt.Errorf("cannot create directory %q: %w", path, err)
	}
	return nil
}

// DoesPostbootExist returns true if the postboot file exists, otherwise false.
func DoesPostbootExist(path string) bool {
	path = fmt.Sprintf("%s/%s", path, constants.PostbootFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Info("Postboot file does not exist", "path", path)
		return false
	}
	return true
}

// GetCmdPath returns the given executable from the PATH, otherwise an error.
func GetCmdPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}
