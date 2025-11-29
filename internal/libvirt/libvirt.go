// Package libvirt is an overlay to the actual libvirt library.
package libvirt

import (
	"errors"
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
		slog.Error("Cannot establish connection with libvirt", "err", err)
		return nil, fmt.Errorf("cannot establish connection with libvirt: %w", err)
	}

	return conn, nil
}

// GetDomain retrieves a given Domain from libvirt.
func GetDomain(conn *libvirt.Connect, image *images.Image, fail bool) (*libvirt.Domain, error) {
	dom, err := conn.LookupDomainByName(image.Name)
	if err != nil {
		if fail {
			slog.Error("Domain is down", "image", image.DisplayName)
			return nil, fmt.Errorf("cannot use %s: domain is down", image.DisplayName)
		}
		return nil, nil
	}

	return dom, nil
}

// GetDomainIPAddress retrieves the IP address of the Domain.
func GetDomainIPAddress(dom *libvirt.Domain, image *images.Image) (string, error) {
	ifaces, err := dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
	if err != nil {
		slog.Debug("Cannot retrieve VM IP address", "err", err, "image", image.Name)
		return "", fmt.Errorf("cannot retrieve VM IP address: %w", err)
	}

	for _, iface := range ifaces {
		if iface.Hwaddr == image.MacAddress {
			return iface.Addrs[0].Addr, nil
		}
	}

	slog.Debug("Cannot retrieve VM IP address", "image", image.Name)
	return "", errors.New("cannot retrieve VM IP address")
}

// ResolveDomainState translates the Domain status into a readable format.
func ResolveDomainState(state libvirt.DomainState) (string, error) {
	display, ok := domainStates[state]
	if !ok {
		slog.Error("Cannot resolve domain state", "state", state)
		return "", fmt.Errorf("cannot resolve domain state: %d", state)
	}

	return display, nil
}
