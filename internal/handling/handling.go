package handling

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
	rawLibvirt "libvirt.org/libvirt-go"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Warn(err)
	}
}

func CloseConnect(c *rawLibvirt.Connect) {
	if _, err := c.Close(); err != nil {
		log.Warn(err)
	}
}

func FreeDomain(d *rawLibvirt.Domain) {
	if err := d.Free(); err != nil {
		log.Warn(err)
	}
}

func ReleaseProcess(p *os.Process) {
	if err := p.Release(); err != nil {
		log.Warn(err)
	}
}
