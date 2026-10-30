package commands

import (
	"fmt"
	"log/slog"

	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/guest"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
)

// StatusCommand represents the options specific to the status command.
type StatusCommand struct{}

// Run is the function for the status command.
func (*StatusCommand) Run(_ *options.Options) error {
	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("Failed to connect to libvirt for status command", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	for _, image := range images.GetAllImages() {
		dom, err := libvirt.GetDomain(conn, &image, false)
		if err != nil {
			slog.Error("Failed to lookup domain for status command", "image", image.Name, "err", err)
			return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
		}

		state, address := "DOWN", "-"
		if dom != nil {
			state, address, err = describeDomain(dom, &image)
			handling.FreeDomain(dom)
			if err != nil {
				return err
			}
		}

		fmt.Printf("%s\t%s\t%s\n", image.Name, state, address)
	}

	return nil
}

func describeDomain(dom *rawLibvirt.Domain, image *images.Image) (string, string, error) {
	info, err := dom.GetInfo()
	if err != nil {
		slog.Error("failed to get domain info", "image", image.Name, "err", err)
		return "", "", fmt.Errorf("cannot get domain info for %q: %w", image.Name, err)
	}

	state, err := libvirt.ResolveDomainState(info.State)
	if err != nil {
		slog.Error("failed to resolve domain state", "image", image.Name, "err", err)
		return "", "", fmt.Errorf("cannot resolve domain state for %q: %w", image.Name, err)
	}

	port, err := guest.SSHPort(dom)
	if err != nil {
		slog.Warn("failed to resolve ssh port", "image", image.Name, "err", err)
		return state, "-", nil
	}

	return state, guest.Address(port), nil
}
