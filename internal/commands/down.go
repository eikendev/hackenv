package commands

import (
	"log"

	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/settings"
)

type DownCommand struct {
}

func (c *DownCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *DownCommand) Run(s *settings.Settings) {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer conn.Close()

	dom := libvirt.GetDomain(conn, &image, true)
	defer dom.Free()

	err := dom.Destroy()
	if err != nil {
		log.Fatalf("Cannot destroy domain: %s\n", err)
	}
}
