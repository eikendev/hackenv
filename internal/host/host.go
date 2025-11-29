// Package host provides various utilities related to the host.
package host

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/paths"
)

// GetHostIPAddress retrieves the IP address of the host on the given interface.
func GetHostIPAddress(ifaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		slog.Error("Cannot retrieve host interfaces", "err", err)
		return "", fmt.Errorf("cannot retrieve host interfaces: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Name != ifaceName {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			ip = ip.To16()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	slog.Error("Cannot retrieve host IP address", "interface", ifaceName)
	return "", fmt.Errorf("cannot retrieve host IP address for interface %s", ifaceName)
}

// GetHostKeyboardLayout retrieves the configured keyboard layout on the host.
func GetHostKeyboardLayout() (string, error) {
	setxkbmapPath, err := paths.GetCmdPath("setxkbmap")
	if err != nil {
		slog.Error("setxkbmap command not found", "err", err)
		return "", fmt.Errorf("cannot locate setxkbmap command: %w", err)
	}

	out, err := exec.Command(
		setxkbmapPath,
		"-query",
	).Output() //#nosec G204
	if err != nil {
		slog.Error("Failed to query keyboard layout", "err", err)
		return "", fmt.Errorf("failed to query keyboard layout: %w", err)
	}

	var line string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "layout") {
			break
		}
	}

	if line == "" {
		slog.Error("Unable to retrieve host keyboard layout: layout line missing")
		return "", fmt.Errorf("unable to retrieve host keyboard layout: layout line missing")
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		slog.Error("Unable to retrieve host keyboard layout: malformed output", "output", line)
		return "", fmt.Errorf("unable to retrieve host keyboard layout: malformed output")
	}
	return parts[len(parts)-1], nil
}
