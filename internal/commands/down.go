package commands

import (
	"fmt"
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
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("Failed to get image details for down command", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("Failed to connect to libvirt for down command", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	dom, err := libvirt.GetDomain(conn, &image, true)
	if err != nil {
		slog.Error("Failed to lookup domain for down command", "image", image.Name, "err", err)
		return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
	}
	if dom != nil {
		defer handling.FreeDomain(dom)
	}

	if err := dom.Destroy(); err != nil {
		slog.Error("Cannot destroy domain", "err", err)
		return fmt.Errorf("cannot destroy domain %q: %w", image.Name, err)
	}

	return nil
}
