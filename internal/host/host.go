// Package host provides various utilities related to the host.
package host

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/paths"
)

type layoutSource struct {
	command string
	args    []string
	key     string
}

var layoutSources = []layoutSource{
	{command: "setxkbmap", args: []string{"-query"}, key: "layout"},
	{command: "localectl", args: []string{"status"}, key: "x11 layout"},
}

// GetHostKeyboardLayout retrieves the configured keyboard layout on the host.
func GetHostKeyboardLayout() (string, error) {
	var errs []error

	for _, source := range layoutSources {
		layout, err := queryLayout(source)
		if err == nil {
			return layout, nil
		}

		slog.Warn("failed to query keyboard layout", "command", source.command, "err", err)
		errs = append(errs, err)
	}

	return "", fmt.Errorf("cannot determine keyboard layout: %w", errors.Join(errs...))
}

func queryLayout(source layoutSource) (string, error) {
	path, err := paths.GetCmdPath(source.command)
	if err != nil {
		return "", fmt.Errorf("cannot locate %s: %w", source.command, err)
	}

	out, err := exec.Command(path, source.args...).Output() //#nosec G204
	if err != nil {
		return "", fmt.Errorf("%s failed: %w", source.command, err)
	}

	return parseLayout(string(out), source.key)
}

func parseLayout(output, key string) (string, error) {
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		name, value, found := strings.Cut(scanner.Text(), ":")
		if !found || !strings.EqualFold(strings.TrimSpace(name), key) {
			continue
		}

		value = strings.TrimSpace(value)
		if value == "" {
			return "", fmt.Errorf("found empty value for %q", key)
		}

		return value, nil
	}

	return "", fmt.Errorf("missing %q in output", key)
}
