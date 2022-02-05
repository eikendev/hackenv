package commands

import (
	"fmt"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/paths"
	"github.com/eikendev/hackenv/internal/settings"
)

type SSHCommand struct{}

func (c *SSHCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *SSHCommand) Run(s *settings.Settings) {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer conn.Close()

	dom := libvirt.GetDomain(conn, &image, true)
	defer dom.Free()

	ipAddr, err := libvirt.GetDomainIPAddress(dom, &image)
	if err != nil {
		log.Fatalf("Cannot retrieve guest's IP address\n")
	}

	args := []string{
		paths.GetCmdPathOrExit("ssh"),
		"-i", paths.GetDataFilePath(constants.SSHKeypairName),
		"-S", "none",
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
		"-X",
		fmt.Sprintf("%s@%s", image.SSHUser, ipAddr),
	}

	if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
		log.Printf("Cannot spawn process: %s\n", err)
	}
}
