// Package host provides various utilities related to the host.
package host

import (
	"bufio"
	"log/slog"
	"net"
	"os"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/paths"
)

// GetHostIPAddress retrieves the IP address of the host on the given interface.
func GetHostIPAddress(ifaceName string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		slog.Error("Cannot retrieve host interfaces", "err", err)
		os.Exit(1)
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

			return ip.String()
		}
	}

	slog.Error("Cannot retrieve host IP address", "interface", ifaceName)
	os.Exit(1)
	return "" // Does not actually return.
}

// GetHostKeyboardLayout retrieves the configured keyboard layout on the host.
func GetHostKeyboardLayout() string {
	out, err := exec.Command(
		paths.GetCmdPathOrExit("setxkbmap"),
		"-query",
	).Output() //#nosec G204
	if err != nil {
		slog.Error("Failed to query keyboard layout", "err", err)
		os.Exit(1)
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
		os.Exit(1)
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	if len(parts) < 2 {
		slog.Error("Unable to retrieve host keyboard layout: malformed output", "output", line)
		os.Exit(1)
		return ""
	}
	return parts[len(parts)-1]
}
