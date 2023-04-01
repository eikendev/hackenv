// Package handling provides convenience functions for cleaning up resources.
package handling

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	rawLibvirt "libvirt.org/libvirt-go"
)

// Close closes an io resource and prints a warning if that fails.
func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Warn(err)
	}
}

// CloseConnect closes a libvirt connections and prints a warning if that fails.
func CloseConnect(c *rawLibvirt.Connect) {
	if _, err := c.Close(); err != nil {
		log.Warn(err)
	}
}

// FreeDomain frees a libvirt domain and prints a warning if that fails.
func FreeDomain(d *rawLibvirt.Domain) {
	if err := d.Free(); err != nil {
		log.Warn(err)
	}
}

// ReleaseProcess releases process information and prints a warning if that fails.
func ReleaseProcess(p *os.Process) {
	if err := p.Release(); err != nil {
		log.Warn(err)
	}
}
