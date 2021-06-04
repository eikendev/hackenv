package commands

import (
	"log"
	"os"
	"os/exec"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/settings"
)

const (
	binPath = "virt-viewer"
)

type GuiCommand struct {
}

func (c *GuiCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *GuiCommand) Run(s *settings.Settings) {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer conn.Close()

	// Check if the domain is up.
	dom := libvirt.GetDomain(conn, &image, true)
	defer dom.Free()

	args := []string{
		binPath,
		"--connect",
		constants.ConnectURI,
		image.Name,
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Printf("Cannot get current working directory: %s\n", err)
	}

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = cwd
	cmd.Env = os.Environ()

	err = cmd.Start()
	if err != nil {
		log.Printf("Cannot spawn process: %s\n", err)
	}
	defer cmd.Process.Release()
}
