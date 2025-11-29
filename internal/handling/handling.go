// Package handling provides convenience functions for cleaning up resources.
package handling

import (
	"io"
	"log/slog"
	"os"

	rawLibvirt "libvirt.org/go/libvirt"
)

// Close closes an io resource and prints a warning if that fails.
func Close(c io.Closer) {
	if c == nil {
		return
	}
	if err := c.Close(); err != nil {
		slog.Warn("Failed to close resource", "err", err)
	}
}

// CloseConnect closes a libvirt connections and prints a warning if that fails.
func CloseConnect(c *rawLibvirt.Connect) {
	if c == nil {
		return
	}
	if _, err := c.Close(); err != nil {
		slog.Warn("Failed to close libvirt connection", "err", err)
	}
}

// FreeDomain frees a libvirt domain and prints a warning if that fails.
func FreeDomain(d *rawLibvirt.Domain) {
	if d == nil {
		return
	}
	if err := d.Free(); err != nil {
		slog.Warn("Failed to free libvirt domain", "err", err)
	}
}

// ReleaseProcess releases process information and prints a warning if that fails.
func ReleaseProcess(p *os.Process) {
	if p == nil {
		return
	}
	if err := p.Release(); err != nil {
		slog.Warn("Failed to release process", "err", err)
	}
}
