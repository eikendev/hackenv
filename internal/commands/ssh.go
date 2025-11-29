package commands

import (
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/eikendev/hackenv/internal/constants"
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
		slog.Error("Failed to get image details for ssh command", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("Failed to connect to libvirt for ssh command", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	dom, err := libvirt.GetDomain(conn, &image, true)
	if err != nil {
		slog.Error("Failed to lookup domain for ssh command", "image", image.Name, "err", err)
		return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
	}
	defer handling.FreeDomain(dom)

	ipAddr, err := libvirt.GetDomainIPAddress(dom, &image)
	if err != nil {
		slog.Error("Cannot retrieve guest IP address", "err", err)
		return fmt.Errorf("cannot retrieve guest IP address: %w", err)
	}

	args, err := buildSSHArgs([]string{
		"-X",
		fmt.Sprintf("%s@%s", image.SSHUser, ipAddr),
	})
	if err != nil {
		slog.Error("Failed to build SSH arguments", "err", err)
		return fmt.Errorf("cannot build SSH arguments: %w", err)
	}

	//#nosec G204
	err = syscall.Exec(args[0], args, os.Environ())
	if err != nil {
		slog.Error("Cannot spawn SSH process", "err", err)
		return fmt.Errorf("failed to exec SSH client: %w", err)
	}

	return nil
}

func buildSSHArgs(customArgs []string) ([]string, error) {
	sshPath, err := paths.GetCmdPath("ssh")
	if err != nil {
		slog.Error("ssh command not found", "err", err)
		return nil, fmt.Errorf("cannot locate ssh binary: %w", err)
	}
	keyPath, err := paths.GetDataFilePath(constants.SSHKeypairName)
	if err != nil {
		slog.Error("Failed to locate SSH keypair path", "err", err)
		return nil, fmt.Errorf("cannot locate ssh keypair: %w", err)
	}

	args := []string{
		sshPath,
		"-i", keyPath,
		"-S", "none",
		"-o", "LogLevel=ERROR",
		"-o", "StrictHostKeyChecking=no",
		"-o", "UserKnownHostsFile=/dev/null",
	}
	args = append(args, customArgs...)

	return args, nil
}
