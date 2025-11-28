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
func GetDataFilePath(file string) string {
	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, file))
	if err != nil {
		slog.Error("Cannot access data directory", "file", file, "err", err)
		os.Exit(1)
	}

	return path
}

// EnsureDirExists creates the given directory if it does not exists.
func EnsureDirExists(path string) {
	err := os.MkdirAll(path, 0o750)
	if err != nil {
		slog.Error("Cannot create directory", "path", path, "err", err)
		os.Exit(1)
	}
}

// DoesPostbootExist returns true if the postboot file exists, otherwise false.
func DoesPostbootExist(path string) bool {
	path = fmt.Sprintf("%s/%s", path, constants.PostbootFile)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		slog.Warn("Postboot file does not exist", "path", path)
		return false
	}
	return true
}

// GetCmdPathOrExit returns the given executable from the PATH.
func GetCmdPathOrExit(cmd string) string {
	path, err := exec.LookPath(cmd)
	if err != nil {
		slog.Error("Command not found", "command", cmd)
		os.Exit(1)
	}

	return path
}

// GetCmdPath returns the given executable from the PATH, otherwise an error.
func GetCmdPath(cmd string) (string, error) {
	return exec.LookPath(cmd)
}
