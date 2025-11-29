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
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("Failed to get image details for GUI command", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("Failed to connect to libvirt for GUI command", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	// Check if the domain is up.
	dom, err := libvirt.GetDomain(conn, &image, true)
	if err != nil {
		slog.Error("Failed to lookup domain for GUI command", "image", image.Name, "err", err)
		return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
	}
	defer handling.FreeDomain(dom)

	args, err := c.viewerArgs(&image)
	if err != nil {
		slog.Error("Failed to resolve viewer arguments", "viewer", c.Viewer, "err", err)
		return fmt.Errorf("cannot resolve viewer %q: %w", c.Viewer, err)
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
		return fmt.Errorf("cannot start viewer %q: %w", args[0], err)
	}
	defer handling.ReleaseProcess(cmd.Process)

	return nil
}

func (c *GuiCommand) viewerArgs(image *images.Image) ([]string, error) {
	switch c.Viewer {
	case virtViewerBin:
		virtViewerPath, err := paths.GetCmdPath(virtViewerBin)
		if err != nil {
			slog.Error("Unable to locate viewer", "viewer", c.Viewer, "err", err)
			return nil, fmt.Errorf("unable to locate %s", c.Viewer)
		}

		args := []string{
			virtViewerPath,
			"--connect",
			constants.ConnectURI,
			image.Name,
		}
		if c.Fullscreen {
			args = append(args, "--full-screen")
		}
		return args, nil
	case remminaBin:
		remminaPath, err := paths.GetCmdPath(remminaBin)
		if err != nil {
			slog.Error("Unable to locate viewer", "viewer", c.Viewer, "err", err)
			return nil, fmt.Errorf("unable to locate %s", c.Viewer)
		}

		return []string{
			remminaPath,
			"-c",
			"SPICE://localhost",
		}, nil
	default:
		slog.Error("Unsupported viewer", "viewer", c.Viewer)
		return nil, fmt.Errorf("cannot use viewer %s: unsupported", c.Viewer)
	}
}
