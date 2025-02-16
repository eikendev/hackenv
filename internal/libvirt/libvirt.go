// Package libvirt is an overlay to the actual libvirt library.
package libvirt

import (
	"errors"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
	"libvirt.org/go/libvirt"
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
		log.Fatalf("Cannot establish connection with libvirt: %s\n", err)
	}

	return conn
}

// GetDomain retrieves a given Domain from libvirt.
func GetDomain(conn *libvirt.Connect, image *images.Image, fail bool) *libvirt.Domain {
	dom, err := conn.LookupDomainByName(image.Name)
	if err != nil {
		if fail {
			log.Fatalf("%s is down\n", image.DisplayName)
		}
		dom = nil
	}

	return dom
}

// GetDomainIPAddress retrieves the IP address of the Domain.
func GetDomainIPAddress(dom *libvirt.Domain, image *images.Image) (string, error) {
	ifaces, err := dom.ListAllInterfaceAddresses(libvirt.DOMAIN_INTERFACE_ADDRESSES_SRC_ARP)
	if err != nil {
		log.Fatalf("Cannot retrieve VM's IP address: %s\n", err)
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
		log.Fatalf("Cannot resolve domain state: %d\n", state)
	}

	return display
}
