package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/settings"
)

type StatusCommand struct{}

func (c *StatusCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *StatusCommand) Run(s *settings.Settings) {
	conn := libvirt.Connect()
	defer conn.Close()

	for _, image := range images.GetAllImages() {
		var state string

		dom := libvirt.GetDomain(conn, &image, false)
		if dom == nil {
			state = "DOWN"
		} else {
			defer dom.Free()

			if info, err := dom.GetInfo(); err != nil {
				log.Printf("Cannot get domain info: %s\n", err)
				continue
			} else {
				state = libvirt.ResolveDomainState(info.State)
			}
		}

		fmt.Printf("%s\t%s\n", image.Name, state)
	}
}
