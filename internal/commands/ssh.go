package commands

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"syscall"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/guest"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/paths"
)

// SSHCommand represents the options specific to the ssh command.
type SSHCommand struct{}

// Run is the function for the ssh command.
func (c *SSHCommand) Run(s *options.Options) error {
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("failed to resolve image", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("failed to connect to libvirt", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	dom, err := libvirt.GetDomain(conn, &image, true)
	if err != nil {
		slog.Error("failed to look up domain", "image", image.Name, "err", err)
		return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
	}
	defer handling.FreeDomain(dom)

	port, err := guest.SSHPort(dom)
	if err != nil {
		slog.Error("failed to resolve ssh port", "image", image.Name, "err", err)
		return fmt.Errorf("cannot determine SSH port of %q: %w", image.Name, err)
	}

	args, err := buildSSHArgs(port, []string{
		"-X",
		fmt.Sprintf("%s@%s", image.SSHUser, guest.Host),
	})
	if err != nil {
		slog.Error("failed to build ssh arguments", "err", err)
		return fmt.Errorf("cannot build SSH arguments: %w", err)
	}

	//#nosec G204
	err = syscall.Exec(args[0], args, os.Environ())
	if err != nil {
		slog.Error("failed to exec ssh", "err", err)
		return fmt.Errorf("failed to exec SSH client: %w", err)
	}

	return nil
}

func buildSSHArgs(port int, customArgs []string) ([]string, error) {
	sshPath, err := paths.GetCmdPath("ssh")
	if err != nil {
		slog.Error("missing ssh command", "err", err)
		return nil, fmt.Errorf("cannot locate ssh binary: %w", err)
	}
	keyPath, err := paths.GetDataFilePath(constants.SSHKeypairName)
	if err != nil {
		slog.Error("failed to resolve ssh keypair path", "err", err)
		return nil, fmt.Errorf("cannot locate ssh keypair: %w", err)
	}

	args := []string{
		sshPath,
		"-i", keyPath,
		"-p", strconv.Itoa(port),
		"-S", "none",
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	args = append(args, customArgs...)

	return args, nil
}
