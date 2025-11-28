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
	conn := libvirt.Connect()
	defer handling.CloseConnect(conn)

	for _, image := range images.GetAllImages() {
		var state string

		image := image
		dom := libvirt.GetDomain(conn, &image, false)
		if dom == nil {
			state = "DOWN"
		} else {
			defer handling.FreeDomain(dom)

			info, err := dom.GetInfo()
			if err != nil {
				slog.Warn("Cannot get domain info", "err", err, "domain", image.Name)
				continue
			}

			state = libvirt.ResolveDomainState(info.State)
		}

		fmt.Printf("%s\t%s\n", image.Name, state)
	}

	return nil
}
