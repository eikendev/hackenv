package commands

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
)

type StatusCommand struct{}

func (c *StatusCommand) Run(s *options.Options) error {
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

			if info, err := dom.GetInfo(); err != nil {
				log.Printf("Cannot get domain info: %s\n", err)
				continue
			} else {
				state = libvirt.ResolveDomainState(info.State)
			}
		}

		fmt.Printf("%s\t%s\n", image.Name, state)
	}

	return nil
}
