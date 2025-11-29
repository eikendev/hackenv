package commands

import (
	"fmt"
	"log/slog"

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
		var state string

		image := image
		dom, err := libvirt.GetDomain(conn, &image, false)
		if err != nil {
			slog.Error("Failed to lookup domain for status command", "image", image.Name, "err", err)
			return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
		}
		if dom == nil {
			state = "DOWN"
		} else {
			defer handling.FreeDomain(dom)

			info, err := dom.GetInfo()
			if err != nil {
				slog.Warn("Cannot get domain info", "err", err, "domain", image.Name)
				continue
			}

			state, err = libvirt.ResolveDomainState(info.State)
			if err != nil {
				slog.Error("Failed to resolve domain state", "image", image.Name, "err", err)
				return fmt.Errorf("cannot resolve domain state for %q: %w", image.Name, err)
			}
		}

		fmt.Printf("%s\t%s\n", image.Name, state)
	}

	return nil
}
