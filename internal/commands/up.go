package commands

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/banner"
	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/guest"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/host"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/paths"
)

const (
	sharedDir        = "shared"
	maxPort          = 65535
	bootGraceDelay   = 20 * time.Second
	sshReadyTimeout  = 2 * time.Minute
	sshProbeWindow   = 10 * time.Second
	sshProbeInterval = time.Second
	xmlTemplate      = `
    <domain type='kvm'>
      <name>%s</name>
      <memory unit='KiB'>%d</memory>
      <vcpu placement='static'>%d</vcpu>
      <os>
        <type>hvm</type>
        <boot dev='cdrom'/>
      </os>
      <features>
        <acpi/>
        <apic/>
      </features>
      <devices>
        <disk type='file' device='cdrom'>
          <driver name='qemu' type='raw'/>
          <source file='%s'/>
          <target dev='sda' bus='sata'/>
          <readonly/>
        </disk>
        <filesystem type='mount' accessmode='mapped'>
          <source dir='%s'/>
          <target dir='/shared'/>
        </filesystem>
        <interface type='user'>
          <mac address='%s'/>
          <model type='virtio'/>
          <backend type='passt'/>
          <portForward proto='tcp' address='%s'>
            <range start='%d' to='22'/>
          </portForward>
        </interface>
        <console type='pty'>
          <target type='serial'/>
        </console>
        <graphics type='spice' port='-1' autoport='yes'>
          <listen type='address' address='127.0.0.1'/>
          <image compression='off'/>
        </graphics>
        <channel type="spicevmc">
            <target type="virtio" name="com.redhat.spice.0"/>
        </channel>
        <video>
          <model type='qxl'/>
        </video>
        <sound model='ich6'/>
        <input type='mouse' bus='ps2'/>
        <input type='keyboard' bus='ps2'/>
        <rng model='virtio'>
          <backend model='random'>/dev/urandom</backend>
        </rng>
      </devices>
    </domain>
`
)

// UpCommand represents the options specific to the up command.
type UpCommand struct {
	Cores       int    `name:"cores" env:"HACKENV_CORES" default:"2" help:"How many virtual CPU cores to assign to the VM"`
	Memory      int    `name:"memory" env:"HACKENV_MEMORY" default:"2097152" help:"How much RAM to assign to the VM (KiB)"`
	SSHPort     int    `name:"ssh-port" env:"HACKENV_SSH_PORT" default:"0" help:"The host port forwarded to the SSH server of the VM (default depends on the image)"`
	DisplaySize string `name:"display_size" env:"HACKENV_DISPLAY_SIZE" default:"1920x1080" help:"The resolution of the VM's display"`
}

func (c *UpCommand) resolveSSHPort(image *images.Image) (int, error) {
	if c.SSHPort == 0 {
		return image.SSHPort, nil
	}

	if c.SSHPort < 1 || c.SSHPort > maxPort {
		slog.Error("found invalid ssh port", "port", c.SSHPort)
		return 0, fmt.Errorf("invalid SSH port %d: must be between 1 and %d", c.SSHPort, maxPort)
	}

	return c.SSHPort, nil
}

func buildXML(c *UpCommand, image *images.Image, path string, port int) (string, error) {
	sharedPath, err := paths.GetDataFilePath(sharedDir)
	if err != nil {
		slog.Error("failed to resolve shared directory", "err", err)
		return "", fmt.Errorf("cannot resolve shared directory path: %w", err)
	}

	if err := paths.EnsureDirExists(sharedPath); err != nil {
		slog.Error("failed to create shared directory", "path", sharedPath, "err", err)
		return "", fmt.Errorf("cannot ensure shared directory at %s: %w", sharedPath, err)
	}

	return fmt.Sprintf(
		xmlTemplate,
		image.Name,
		c.Memory,
		c.Cores,
		path,
		sharedPath,
		image.MacAddress,
		guest.Host,
		port,
	), nil
}

func pollSSH(port int, window time.Duration) bool {
	timer := time.NewTimer(window)
	defer timer.Stop()

	ticker := time.NewTicker(sshProbeInterval)
	defer ticker.Stop()

	for {
		if guest.IsSSHReachable(port, sshProbeInterval) {
			return true
		}

		select {
		case <-timer.C:
			return false
		case <-ticker.C:
		}
	}
}

func waitForSSH(dom *rawLibvirt.Domain, image *images.Image, port int) error {
	slog.Info("waiting for guest boot", "image", image.Name)
	time.Sleep(bootGraceDelay)

	deadline := time.Now().Add(sshReadyTimeout)
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		slog.Info("starting ssh in guest", "attempt", attempt, "image", image.Name)

		if err := image.StartSSH(dom); err != nil {
			slog.Error("failed to start ssh in guest", "image", image.Name, "err", err)
			return fmt.Errorf("failed to start SSH in %s: %w", image.DisplayName, err)
		}

		if pollSSH(port, sshProbeWindow) {
			slog.Info("reached guest ssh", "address", guest.Address(port), "image", image.Name)
			return nil
		}
	}

	slog.Error("failed to reach guest ssh", "image", image.Name, "address", guest.Address(port), "timeout", sshReadyTimeout)
	return fmt.Errorf("SSH on %s did not become reachable within %s; run 'hackenv doctor'", guest.Address(port), sshReadyTimeout)
}

func provisionClient(image *images.Image, port int) error {
	sharedPath, err := paths.GetDataFilePath(sharedDir)
	if err != nil {
		slog.Error("failed to resolve shared directory", "err", err)
		return fmt.Errorf("cannot locate shared directory: %w", err)
	}

	if !paths.DoesPostbootExist(sharedPath) {
		return nil
	}

	args, err := buildSSHArgs(port, []string{
		fmt.Sprintf("%s@%s", image.SSHUser, guest.Host),
		fmt.Sprintf("/shared/%s", constants.PostbootFile),
	})
	if err != nil {
		slog.Error("failed to build ssh arguments", "err", err)
		return fmt.Errorf("cannot build provisioning SSH args: %w", err)
	}

	slog.Info("provisioning guest", "image", image.Name)

	//#nosec G204
	err = syscall.Exec(args[0], args, os.Environ())
	if err != nil {
		slog.Error("failed to provision guest", "err", err)
		return fmt.Errorf("failed to execute provisioning command: %w", err)
	}

	return nil
}

func connectGuest(image *images.Image, port int) (*goph.Client, error) {
	client, err := goph.NewConn(&goph.Config{
		User:     image.SSHUser,
		Addr:     guest.Host,
		Port:     uint(port), //#nosec G115 -- port is validated to fit in 16 bits
		Auth:     goph.Password(image.SSHPassword),
		Timeout:  goph.DefaultTimeout,
		Callback: ssh.InsecureIgnoreHostKey(), //#nosec G106 -- the guest is a fresh live VM on loopback
	})
	if err != nil {
		slog.Error("failed to connect to guest", "image", image.Name, "err", err)
		return nil, fmt.Errorf("cannot create SSH client for %s: %w", image.Name, err)
	}
	if client == nil {
		slog.Error("missing ssh client", "image", image.Name)
		return nil, fmt.Errorf("encountered nil SSH client for %s", image.Name)
	}

	return client, nil
}

func configureClient(c *UpCommand, image *images.Image, port int, keymap string) error {
	client, err := connectGuest(image, port)
	if err != nil {
		return err
	}

	publicKeyPath, err := paths.GetDataFilePath(constants.SSHKeypairName + ".pub")
	if err != nil {
		slog.Error("failed to resolve public key path", "err", err)
		return fmt.Errorf("cannot locate public key: %w", err)
	}

	publicKey, err := os.ReadFile(publicKeyPath) //#nosec G304
	if err != nil {
		slog.Error("failed to read public key", "path", publicKeyPath, "err", err)
		return fmt.Errorf("cannot read public key %s: %w", publicKeyPath, err)
	}

	publicKeyStr := string(publicKey)

	if keymap == "" {
		keymap, err = host.GetHostKeyboardLayout()
		if err != nil {
			slog.Warn("skipping guest keymap configuration", "hint", "use --keymap to set one", "err", err)
		}
	}

	cmds := image.ConfigurationCmds
	cmds = append(cmds, []string{
		// Add the SSH key to authorized_keys.
		"mkdir ~/.ssh",
		"chmod 700 ~/.ssh",
		"printf '" + publicKeyStr + "' >> ~/.ssh/authorized_keys",
		"chmod 660 ~/.ssh/authorized_keys",

		// Disable password authentication on SSH.
		"sudo sed -i '/PasswordAuthentication/s/yes/no/' /etc/ssh/sshd_config",
		"sudo systemctl reload ssh",

		// Setup a shared directory.
		"sudo mkdir /shared",
		"sudo mount -t 9p -o trans=virtio,version=9p2000.L /shared /shared",

		// Set screen size to Full HD.
		fmt.Sprintf("DISPLAY=:0 xrandr --size %s", c.DisplaySize),
	}...)

	if keymap != "" {
		cmds = append(cmds, fmt.Sprintf("DISPLAY=:0 setxkbmap %s", keymap))
	}

	for _, cmd := range cmds {
		_, err := client.Run(cmd)
		if err != nil {
			slog.Error("failed to run guest command", "command", cmd, "err", err)
			return fmt.Errorf("failed to run command '%s' over SSH: %s", cmd, err)
		}
	}

	return nil
}

func ensureSSHKeypairExists() error {
	sshKeypairPath, err := paths.GetDataFilePath(constants.SSHKeypairName)
	if err != nil {
		slog.Error("failed to resolve ssh keypair path", "err", err)
		return fmt.Errorf("cannot determine SSH keypair path: %w", err)
	}

	if _, err := os.Stat(sshKeypairPath); err == nil {
		// SSH keypair already exists.
		return nil
	}

	sshKeygenPath, err := paths.GetCmdPath("ssh-keygen")
	if err != nil {
		slog.Error("missing ssh-keygen command", "err", err)
		return fmt.Errorf("cannot locate ssh-keygen: %w", err)
	}

	cmd := exec.Command(
		sshKeygenPath,
		"-f",
		sshKeypairPath,
		"-t",
		"ed25519",
		"-C",
		constants.SSHKeypairName,
		"-q",
		"-N",
		"", // Password is empty so no typing is required.
	) //#nosec G204

	if err := cmd.Start(); err != nil {
		slog.Error("failed to start ssh-keygen", "err", err)
		return fmt.Errorf("cannot start ssh-keygen: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("failed to run ssh-keygen", "err", err)
		return fmt.Errorf("failed to finish ssh-keygen: %w", err)
	}

	return nil
}

// Run is the function for the up command.
func (c *UpCommand) Run(s *options.Options) error {
	banner.PrintBanner()

	image, localPath, localVersion, err := c.resolveImage(s)
	if err != nil {
		slog.Error("failed to resolve image", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image data: %w", err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("failed to connect to libvirt", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	existing, err := libvirt.GetDomain(conn, &image, false)
	if err != nil {
		return fmt.Errorf("cannot look up domain %q: %w", image.Name, err)
	}
	if existing != nil {
		defer handling.FreeDomain(existing)
		return c.handleExistingDomain(existing, &image, s)
	}

	port, err := c.resolveSSHPort(&image)
	if err != nil {
		return err
	}

	if err := guest.EnsurePortFree(port); err != nil {
		return fmt.Errorf("%w; choose another port with --ssh-port or HACKENV_SSH_PORT", err)
	}

	xml, err := buildXML(c, &image, localPath, port)
	if err != nil {
		slog.Error("failed to build domain xml", "image", image.Name, "err", err)
		return fmt.Errorf("cannot build domain definition: %w", err)
	}

	dom, err := conn.DomainCreateXML(xml, 0)
	if err != nil {
		slog.Error("failed to create domain", "image", image.Name, "hint", "run hackenv doctor", "err", err)
		return fmt.Errorf("cannot create domain %q: %w", image.Name, err)
	}
	defer handling.FreeDomain(dom)

	return c.bootAndConfigure(dom, &image, localVersion, port, s)
}

func (c *UpCommand) resolveImage(s *options.Options) (images.Image, string, string, error) {
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("failed to resolve image", "type", s.Type, "err", err)
		return images.Image{}, "", "", fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	localPath, err := image.GetLatestPath()
	if err != nil {
		slog.Error("failed to resolve image path", "image", image.DisplayName, "err", err)
		return images.Image{}, "", "", fmt.Errorf("cannot resolve latest image path for %s: %w", image.DisplayName, err)
	}

	localVersion := image.FileVersion(localPath)

	if info, err := image.GetDownloadInfo(false); err == nil && info != nil {
		if !image.VersionComparer.Eq(info.Version, localVersion) {
			slog.Info("found new image version", "image", image.DisplayName, "version", info.Version)
		}
	} else if err != nil {
		slog.Warn("failed to fetch download info", "image", image.DisplayName, "err", err)
	}

	if err := ensureSSHKeypairExists(); err != nil {
		slog.Error("failed to create ssh keypair", "err", err)
		return images.Image{}, "", "", fmt.Errorf("failed to ensure SSH keypair exists: %w", err)
	}

	return image, localPath, localVersion, nil
}

func (c *UpCommand) handleExistingDomain(dom *rawLibvirt.Domain, image *images.Image, s *options.Options) error {
	if !s.Provision {
		slog.Error("domain already running", "image", image.Name)
		return fmt.Errorf("cannot create domain %q: already running", image.Name)
	}

	slog.Info("skipping domain creation", "image", image.Name, "reason", "already running")

	port, err := guest.SSHPort(dom)
	if err != nil {
		return fmt.Errorf("cannot determine SSH port of running domain %q: %w", image.Name, err)
	}

	if !pollSSH(port, sshProbeWindow) {
		slog.Error("failed to reach guest ssh", "image", image.Name, "address", guest.Address(port))
		return fmt.Errorf("running domain %q is not reachable on %s", image.Name, guest.Address(port))
	}

	return provisionClient(image, port)
}

func (c *UpCommand) bootAndConfigure(dom *rawLibvirt.Domain, image *images.Image, version string, port int, opts *options.Options) error {
	if err := image.Boot(dom, version); err != nil {
		slog.Error("failed to boot image", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed to boot image %s: %w", image.DisplayName, err)
	}

	if err := waitForSSH(dom, image, port); err != nil {
		return err
	}

	if err := configureClient(c, image, port, opts.Keymap); err != nil {
		slog.Error("failed to configure guest", "err", err)
		return fmt.Errorf("cannot configure guest %s: %w", image.DisplayName, err)
	}

	slog.Info("completed guest setup", "image", image.Name, "address", guest.Address(port))

	if opts.Provision {
		if err := provisionClient(image, port); err != nil {
			return fmt.Errorf("failed to provision guest %s: %w", image.DisplayName, err)
		}
	}

	return nil
}
