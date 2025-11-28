package commands

import (
	"log/slog"

	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
)

// DownCommand represents the options specific to the down command.
type DownCommand struct{}

// Run is the function for the down command.
func (c *DownCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer handling.CloseConnect(conn)

	dom := libvirt.GetDomain(conn, &image, true)
	defer handling.FreeDomain(dom)

	err := dom.Destroy()
	if err != nil {
		slog.Error("Cannot destroy domain", "err", err)
		return err
	}

	return nil
}
