package commands

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/melbahja/goph"
	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/banner"
	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/host"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/paths"
)

const (
	sharedDir    = "shared"
	connectTries = 60
	xmlTemplate  = `
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
        <interface type='bridge'>
          <mac address='%s'/>
          <source bridge='%s'/>
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
	Interface   string `name:"iface" env:"HACKENV_IFACE" default:"virbr0" help:"The network interface to use as a bridge"`
	DisplaySize string `name:"display_size" env:"HACKENV_DISPLAY_SIZE" default:"1920x1080" help:"The resolution of the VM's display"`
}

func buildXML(c *UpCommand, image images.Image, path string) (string, error) {
	sharedPath, err := paths.GetDataFilePath(sharedDir)
	if err != nil {
		slog.Error("Failed to resolve shared directory path", "err", err)
		return "", fmt.Errorf("cannot resolve shared directory path: %w", err)
	}

	if err := paths.EnsureDirExists(sharedPath); err != nil {
		slog.Error("Failed to ensure shared directory exists", "path", sharedPath, "err", err)
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
		c.Interface,
	), nil
}

func waitBootComplete(dom *rawLibvirt.Domain, image *images.Image) (string, error) {
	for i := 1; i <= connectTries; i++ {
		slog.Info("Waiting for VM to become active", "attempt", i, "maxAttempts", connectTries, "image", image.Name)

		ipAddr, err := libvirt.GetDomainIPAddress(dom, image)
		if err == nil {
			slog.Info("VM is up", "ip", ipAddr, "image", image.Name)
			return ipAddr, nil
		}

		time.Sleep(2 * time.Second)
	}

	slog.Error("VM did not become active", "image", image.Name, "attempts", connectTries)
	return "", fmt.Errorf("failed to detect active VM within %d attempts", connectTries)
}

func provisionClient(_ *UpCommand, image *images.Image, guestIPAddr string) error {
	sharedPath, err := paths.GetDataFilePath(sharedDir)
	if err != nil {
		slog.Error("Failed to locate shared path for provisioning", "err", err)
		return fmt.Errorf("cannot locate shared directory: %w", err)
	}

	if paths.DoesPostbootExist(sharedPath) {
		args, err := buildSSHArgs([]string{
			fmt.Sprintf("%s@%s", image.SSHUser, guestIPAddr),
			fmt.Sprintf("/shared/%s", constants.PostbootFile),
		})
		if err != nil {
			slog.Error("Failed to build SSH arguments for provisioning", "err", err)
			return fmt.Errorf("cannot build provisioning SSH args: %w", err)
		}

		slog.Info("Provisioning VM", "image", image.Name)

		//#nosec G204
		err = syscall.Exec(args[0], args, os.Environ())
		if err != nil {
			slog.Error("Cannot provision VM", "err", err)
			return fmt.Errorf("failed to execute provisioning command: %w", err)
		}
	}

	return nil
}

func configureClient(c *UpCommand, _ *rawLibvirt.Domain, image *images.Image, guestIPAddr string, keymap string) error {
	client, err := goph.NewUnknown(image.SSHUser, guestIPAddr, goph.Password(image.SSHPassword))
	if err != nil {
		slog.Error("Failed to create SSH client for guest configuration", "image", image.Name, "err", err)
		return fmt.Errorf("cannot create SSH client for %s: %w", image.Name, err)
	}
	if client == nil {
		slog.Error("SSH client for guest configuration is nil", "image", image.Name)
		return fmt.Errorf("encountered nil SSH client for %s", image.Name)
	}

	publicKeyPath, err := paths.GetDataFilePath(constants.SSHKeypairName + ".pub")
	if err != nil {
		slog.Error("Failed to locate public key for provisioning", "err", err)
		return fmt.Errorf("cannot locate public key: %w", err)
	}

	publicKey, err := os.ReadFile(publicKeyPath) //#nosec G304
	if err != nil {
		slog.Error("Failed to read public key for provisioning", "path", publicKeyPath, "err", err)
		return fmt.Errorf("cannot read public key %s: %w", publicKeyPath, err)
	}

	publicKeyStr := string(publicKey)

	if keymap == "" {
		keymap, err = host.GetHostKeyboardLayout()
		if err != nil {
			slog.Error("Failed to detect host keyboard layout", "err", err)
			return fmt.Errorf("cannot detect host keyboard layout: %w", err)
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

		// Set keyboard layout.
		fmt.Sprintf("DISPLAY=:0 setxkbmap %s", keymap),
	}...)

	for _, cmd := range cmds {
		_, err := client.Run(cmd)
		if err != nil {
			slog.Error("Failed to run guest configuration command", "command", cmd, "err", err)
			return fmt.Errorf("failed to run command '%s' over SSH: %s", cmd, err)
		}
	}

	return nil
}

func ensureSSHKeypairExists() error {
	sshKeypairPath, err := paths.GetDataFilePath(constants.SSHKeypairName)
	if err != nil {
		slog.Error("Failed to determine SSH keypair path", "err", err)
		return fmt.Errorf("cannot determine SSH keypair path: %w", err)
	}

	if _, err := os.Stat(sshKeypairPath); err == nil {
		// SSH keypair already exists.
		return nil
	}

	sshKeygenPath, err := paths.GetCmdPath("ssh-keygen")
	if err != nil {
		slog.Error("ssh-keygen command not found", "err", err)
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
		slog.Error("Failed to start ssh-keygen", "err", err)
		return fmt.Errorf("cannot start ssh-keygen: %w", err)
	}

	if err := cmd.Wait(); err != nil {
		slog.Error("ssh-keygen exited with error", "err", err)
		return fmt.Errorf("failed to finish ssh-keygen: %w", err)
	}

	return nil
}

// Run is the function for the up command.
func (c *UpCommand) Run(s *options.Options) error {
	banner.PrintBanner()

	image, localPath, localVersion, err := c.resolveImage(s)
	if err != nil {
		slog.Error("Failed to resolve image for up command", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image data: %w", err)
	}

	xml, err := buildXML(c, image, localPath)
	if err != nil {
		slog.Error("Failed to build domain XML", "image", image.Name, "err", err)
		return fmt.Errorf("cannot build domain definition: %w", err)
	}

	conn, err := libvirt.Connect()
	if err != nil {
		slog.Error("Failed to connect to libvirt for up command", "err", err)
		return fmt.Errorf("cannot connect to libvirt: %w", err)
	}
	defer handling.CloseConnect(conn)

	dom, err := conn.DomainCreateXML(xml, 0)
	if dom != nil {
		defer handling.FreeDomain(dom)
	}

	if err != nil {
		if s.Provision {
			return c.provisionExistingDomain(conn, &image)
		}

		slog.Error("Cannot create domain. Try running 'hackenv fix all'.", "err", err)
		return fmt.Errorf("cannot create domain %q: %w", image.Name, err)
	}

	return c.bootAndConfigure(dom, &image, localVersion, s)
}

func (c *UpCommand) resolveImage(s *options.Options) (images.Image, string, string, error) {
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("Failed to get image details for up command", "type", s.Type, "err", err)
		return images.Image{}, "", "", fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}

	localPath, err := image.GetLatestPath()
	if err != nil {
		slog.Error("Failed to resolve latest image path", "image", image.DisplayName, "err", err)
		return images.Image{}, "", "", fmt.Errorf("cannot resolve latest image path for %s: %w", image.DisplayName, err)
	}

	localVersion := image.FileVersion(localPath)

	if info, err := image.GetDownloadInfo(false); err == nil && info != nil {
		if !image.VersionComparer.Eq(info.Version, localVersion) {
			slog.Info("New image version available", "image", image.DisplayName, "version", info.Version)
		}
	} else if err != nil {
		slog.Warn("Unable to determine latest upstream image version", "image", image.DisplayName, "err", err)
	}

	if err := ensureSSHKeypairExists(); err != nil {
		slog.Error("Cannot create SSH keypair", "err", err)
		return images.Image{}, "", "", fmt.Errorf("failed to ensure SSH keypair exists: %w", err)
	}

	return image, localPath, localVersion, nil
}

func (c *UpCommand) provisionExistingDomain(conn *rawLibvirt.Connect, image *images.Image) error {
	slog.Info("Domain already running, provisioning instead", "image", image.DisplayName)

	dom, err := libvirt.GetDomain(conn, image, true)
	if err != nil {
		slog.Error("Failed to lookup existing domain for provisioning", "image", image.DisplayName, "err", err)
		return fmt.Errorf("cannot look up running domain %q: %w", image.DisplayName, err)
	}
	defer handling.FreeDomain(dom)

	guestIPAddr, err := waitBootComplete(dom, image)
	if err != nil {
		slog.Error("Existing domain did not become ready for provisioning", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed while waiting for running domain %q: %w", image.DisplayName, err)
	}

	if err := provisionClient(c, image, guestIPAddr); err != nil {
		slog.Error("Provisioning of existing domain failed", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed to provision running domain %q: %w", image.DisplayName, err)
	}

	return nil
}

func (c *UpCommand) bootAndConfigure(dom *rawLibvirt.Domain, image *images.Image, version string, opts *options.Options) error {
	if err := image.Boot(dom, version); err != nil {
		slog.Error("Failed during image boot sequence", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed to boot image %s: %w", image.DisplayName, err)
	}

	guestIPAddr, err := waitBootComplete(dom, image)
	if err != nil {
		slog.Error("Domain did not become ready after boot", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed while waiting for domain %s: %w", image.DisplayName, err)
	}

	if err := image.StartSSH(dom); err != nil {
		slog.Error("Failed to start SSH on guest", "image", image.DisplayName, "err", err)
		return fmt.Errorf("failed to start SSH on %s: %w", image.DisplayName, err)
	}

	if err := configureClient(c, dom, image, guestIPAddr, opts.Keymap); err != nil {
		slog.Error("Cannot configure client", "err", err)
		return fmt.Errorf("cannot configure guest %s: %w", image.DisplayName, err)
	}

	slog.Info("VM is ready to use", "image", image.DisplayName)

	if opts.Provision {
		if err := provisionClient(c, image, guestIPAddr); err != nil {
			return fmt.Errorf("failed to provision guest %s: %w", image.DisplayName, err)
		}
	}

	return nil
}
