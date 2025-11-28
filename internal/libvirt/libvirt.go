// Package libvirt is an overlay to the actual libvirt library.
package libvirt

import (
	"errors"
	"log/slog"
	"os"

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
func Connect() *libvirt.Connect {
	conn, err := libvirt.NewConnect(constants.ConnectURI)
	if err != nil {
		slog.Error("Cannot establish connection with libvirt", "err", err)
		os.Exit(1)
	}

	return conn
}

// GetDomain retrieves a given Domain from libvirt.
func GetDomain(conn *libvirt.Connect, image *images.Image, fail bool) *libvirt.Domain {
	dom, err := conn.LookupDomainByName(image.Name)
	if err != nil {
		if fail {
			slog.Error("Domain is down", "image", image.DisplayName)
			os.Exit(1)
		}
		dom = nil
	}

	return dom
}

// GetDomainIPAddress retrieves the IP address of the Domain.
func GetDomainIPAddress(dom *libvirt.Domain, image *images.Image) (string, error) {
	ifaces, err := dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
	if err != nil {
		slog.Error("Cannot retrieve VM IP address", "err", err, "image", image.Name)
		os.Exit(1)
	}

	for _, iface := range ifaces {
		if iface.Hwaddr == image.MacAddress {
			return iface.Addrs[0].Addr, nil
		}
	}

	return "", errors.New("cannot retrieve VM's IP address")
}

// ResolveDomainState translates the Domain status into a readable format.
func ResolveDomainState(state libvirt.DomainState) string {
	display, ok := domainStates[state]
	if !ok {
		slog.Error("Cannot resolve domain state", "state", state)
		os.Exit(1)
	}

	return display
}
