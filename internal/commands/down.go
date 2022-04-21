package commands

import (
	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
)

type DownCommand struct{}

func (c *DownCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer conn.Close()

	dom := libvirt.GetDomain(conn, &image, true)
	defer dom.Free()

	err := dom.Destroy()
	if err != nil {
		log.Fatalf("Cannot destroy domain: %s\n", err)
	}

	return nil
}
