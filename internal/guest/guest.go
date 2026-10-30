// Package guest provides utilities to reach the guest from the host.
package guest

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strconv"
	"time"

	rawLibvirt "libvirt.org/go/libvirt"
)

// Host is the address on which passt exposes forwarded guest ports.
const Host = "127.0.0.1"

const sshBannerPrefix = "SSH-"

type domainXML struct {
	Devices struct {
		Filesystems []struct {
			Driver struct {
				Type string `xml:"type,attr"`
			} `xml:"driver"`
		} `xml:"filesystem"`
		Interfaces []struct {
			PortForwards []struct {
				Proto  string `xml:"proto,attr"`
				Ranges []struct {
					Start int `xml:"start,attr"`
					To    int `xml:"to,attr"`
				} `xml:"range"`
			} `xml:"portForward"`
		} `xml:"interface"`
	} `xml:"devices"`
}

// SSHPort returns the host port that is forwarded to the guest's SSH port.
func SSHPort(dom *rawLibvirt.Domain) (int, error) {
	desc, err := dom.GetXMLDesc(0)
	if err != nil {
		slog.Error("failed to read domain description", "err", err)
		return 0, fmt.Errorf("cannot read domain description: %w", err)
	}

	return parseSSHPort(desc)
}

func parseSSHPort(desc string) (int, error) {
	var parsed domainXML
	if err := xml.Unmarshal([]byte(desc), &parsed); err != nil {
		slog.Error("failed to parse domain description", "err", err)
		return 0, fmt.Errorf("cannot parse domain description: %w", err)
	}

	for _, iface := range parsed.Devices.Interfaces {
		for _, forward := range iface.PortForwards {
			if forward.Proto != "tcp" {
				continue
			}
			for _, r := range forward.Ranges {
				if r.To == 22 {
					return r.Start, nil
				}
			}
		}
	}

	slog.Error("missing ssh port forward")
	return 0, errors.New("domain has no SSH port forward")
}

// SharedMode returns the driver of the shared directory, "virtiofs" or "9p".
func SharedMode(dom *rawLibvirt.Domain) (string, error) {
	desc, err := dom.GetXMLDesc(0)
	if err != nil {
		slog.Error("failed to read domain description", "err", err)
		return "", fmt.Errorf("cannot read domain description: %w", err)
	}

	return parseSharedMode(desc)
}

func parseSharedMode(desc string) (string, error) {
	var parsed domainXML
	if err := xml.Unmarshal([]byte(desc), &parsed); err != nil {
		slog.Error("failed to parse domain description", "err", err)
		return "", fmt.Errorf("cannot parse domain description: %w", err)
	}

	if len(parsed.Devices.Filesystems) == 0 {
		slog.Error("missing shared filesystem")
		return "", errors.New("domain has no shared filesystem")
	}

	if parsed.Devices.Filesystems[0].Driver.Type == "virtiofs" {
		return "virtiofs", nil
	}

	return "9p", nil
}

// Address returns the host address and port to reach the guest's SSH server.
func Address(port int) string {
	return net.JoinHostPort(Host, strconv.Itoa(port))
}

// EnsurePortFree fails if the given host port cannot be bound.
func EnsurePortFree(port int) error {
	listener, err := net.Listen("tcp", Address(port))
	if err != nil {
		slog.Error("failed to bind ssh port", "port", port, "err", err)
		return fmt.Errorf("port %d is not available on the host: %w", port, err)
	}

	if err := listener.Close(); err != nil {
		slog.Warn("failed to close probe listener", "port", port, "err", err)
	}

	return nil
}

// IsSSHReachable reports whether an SSH server answers on the forwarded port.
// A plain TCP connect is not enough: passt accepts the connection on the host side
// before it knows whether anything listens in the guest, so the SSH banner is required.
func IsSSHReachable(port int, timeout time.Duration) bool {
	conn, err := net.DialTimeout("tcp", Address(port), timeout)
	if err != nil {
		return false
	}
	defer func() {
		if err := conn.Close(); err != nil {
			slog.Warn("failed to close probe connection", "port", port, "err", err)
		}
	}()

	if err := conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return false
	}

	banner := make([]byte, len(sshBannerPrefix))
	if _, err := io.ReadFull(conn, banner); err != nil {
		return false
	}

	return string(banner) == sshBannerPrefix
}
