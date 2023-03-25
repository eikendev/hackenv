package commands

import (
	"fmt"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/paths"
)

type SSHCommand struct{}

func (c *SSHCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer handling.CloseConnect(conn)

	dom := libvirt.GetDomain(conn, &image, true)
	defer handling.FreeDomain(dom)

	ipAddr, err := libvirt.GetDomainIPAddress(dom, &image)
	if err != nil {
		log.Fatalf("Cannot retrieve guest's IP address\n")
	}

	args := buildSSHArgs([]string{
		"-X",
		fmt.Sprintf("%s@%s", image.SSHUser, ipAddr),
	})

	//#nosec G204
	if err := syscall.Exec(args[0], args, os.Environ()); err != nil {
		log.Printf("Cannot spawn process: %s\n", err)
	}

	return nil
}

func buildSSHArgs(customArgs []string) []string {
	args := []string{
		paths.GetCmdPathOrExit("ssh"),
		"-i", paths.GetDataFilePath(constants.SSHKeypairName),
		"-S", "none",
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	args = append(args, customArgs[:]...)

	return args
}
