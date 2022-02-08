package commands

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/paths"
	"github.com/eikendev/hackenv/internal/settings"
)

const (
	virtViewerBin = "virt-viewer"
	remminaBin    = "remmina"
)

type GuiCommand struct {
	//lint:ignore SA5008 go-flags makes use of duplicate struct tags
	Viewer     string `long:"viewer" env:"HACKENV_VIEWER" default:"virt-viewer" choice:"virt-viewer" choice:"remmina" description:"The viewer to use to connect to the VM"`
	Fullscreen bool   `long:"fullscreen" alias:"f" env:"HACKENV_FULLSCREEN" description:"Start GUI in fullscreen (virt-viewer only)"`
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

	var args []string

	if virtViewerPath, err := paths.GetCmdPath(virtViewerBin); c.Viewer == virtViewerBin && err == nil {
		args = []string{
			virtViewerPath,
			"--connect",
			constants.ConnectURI,
			image.Name,
		}

		if c.Fullscreen {
			args = append(args, []string{"--full-screen"}...)
		}

	} else if remminaPath, err := paths.GetCmdPath(remminaBin); c.Viewer == remminaBin && err == nil {
		args = []string{
			remminaPath,
			"-c",
			"SPICE://localhost",
		}
	} else {
		log.Fatalf("Unable to locate %s to connect to the VM.\n", c.Viewer)
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
