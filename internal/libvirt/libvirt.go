// Package libvirt is an overlay to the actual libvirt library.
package libvirt

import (
	"fmt"
	"log/slog"

	"libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
)

var domainStates = map[libvirt.DomainState]string{
	libvirt.DOMAIN_NOSTATE:     "NOSTATE",
	libvirt.DOMAIN_RUNNING:     "RUNNING",
	libvirt.DOMAIN_BLOCKED:     "BLOCKED",
	libvirt.DOMAIN_PAUSED:      "PAUSED",
	libvirt.DOMAIN_SHUTDOWN:    "SHUTDOWN",
	libvirt.DOMAIN_CRASHED:     "CRASHED",
	libvirt.DOMAIN_PMSUSPENDED: "PMSUSPENDED",
	libvirt.DOMAIN_SHUTOFF:     "SHUTOFF",
}

// Connect establishes a connection to libvirt.
func Connect() (*libvirt.Connect, error) {
	conn, err := libvirt.NewConnect(constants.ConnectURI)
	if err != nil {
		slog.Error("failed to connect to libvirt", "err", err)
		return nil, fmt.Errorf("cannot establish connection with libvirt: %w", err)
	}

	return conn, nil
}

// GetDomain retrieves a given Domain from libvirt.
func GetDomain(conn *libvirt.Connect, image *images.Image, fail bool) (*libvirt.Domain, error) {
	dom, err := conn.LookupDomainByName(image.Name)
	if err != nil {
		if fail {
			slog.Error("missing domain", "image", image.DisplayName)
			return nil, fmt.Errorf("cannot use %s: domain is down", image.DisplayName)
		}
		return nil, nil
	}

	return dom, nil
}

// ResolveDomainState translates the Domain status into a readable format.
func ResolveDomainState(state libvirt.DomainState) (string, error) {
	display, ok := domainStates[state]
	if !ok {
		slog.Error("failed to resolve domain state", "state", state)
		return "", fmt.Errorf("cannot resolve domain state: %d", state)
	}

	return display, nil
}
