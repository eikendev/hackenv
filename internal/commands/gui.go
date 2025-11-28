package commands

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/paths"
)

const (
	virtViewerBin = "virt-viewer"
	remminaBin    = "remmina"
)

// GuiCommand represents the options specific to the gui command.
type GuiCommand struct {
	Viewer     string `name:"viewer" env:"HACKENV_VIEWER" default:"virt-viewer" enum:"virt-viewer,remmina" help:"The viewer to use to connect to the VM"`
	Fullscreen bool   `short:"f" name:"fullscreen" env:"HACKENV_FULLSCREEN" help:"Start GUI in fullscreen (virt-viewer only)"`
}

// Run is the function for the gui command.
func (c *GuiCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)

	conn := libvirt.Connect()
	defer handling.CloseConnect(conn)

	// Check if the domain is up.
	dom := libvirt.GetDomain(conn, &image, true)
	defer handling.FreeDomain(dom)

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
		slog.Error("Unable to locate viewer to connect to the VM", "viewer", c.Viewer)
		return fmt.Errorf("unable to locate %s to connect to the VM", c.Viewer)
	}

	cwd, err := os.Getwd()
	if err != nil {
		slog.Warn("Cannot get current working directory", "err", err)
	}

	cmd := exec.Command(args[0], args[1:]...) //#nosec G204
	cmd.Dir = cwd
	cmd.Env = os.Environ()

	err = cmd.Start()
	if err != nil {
		slog.Error("Cannot spawn viewer process", "err", err)
	}
	defer handling.ReleaseProcess(cmd.Process)

	return nil
}
