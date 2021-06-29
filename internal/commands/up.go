package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/host"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/libvirt"
	"github.com/eikendev/hackenv/internal/paths"
	"github.com/eikendev/hackenv/internal/settings"
	"github.com/melbahja/goph"
	rawLibvirt "libvirt.org/libvirt-go"
)

const (
	sharedDir    = "shared"
	connectTries = 20
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

type UpCommand struct {
	Cores       int    `long:"cores" env:"HACKENV_CORES" default:"2" description:"How many virtual CPU cores to assign to the VM"`
	Memory      int    `long:"memory" env:"HACKENV_MEMORY" default:"2097152" description:"How much RAM to assign to the VM"`
	Interface   string `long:"iface" env:"HACKENV_IFACE" default:"virbr0" description:"The network interface to use as a bridge"`
	DisplaySize string `long:"display_size" env:"HACKENV_DISPLAY_SIZE" default:"1920x1080" description:"The resolution of the VM's display"`
}

func (c *UpCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func buildXML(c *UpCommand, image images.Image) string {
	sharedPath := paths.GetDataFilePath(sharedDir)
	paths.EnsureDirExists(sharedPath)

	return fmt.Sprintf(
		xmlTemplate,
		image.Name,
		c.Memory,
		c.Cores,
		image.GetLocalPath(),
		sharedPath,
		image.MacAddress,
		c.Interface,
	)
}

func waitBootComplete(dom *rawLibvirt.Domain, image *images.Image) string {
	for i := 1; i <= connectTries; i++ {
		log.Printf("Waiting for VM to become active (%02d/%d)...\n", i, connectTries)

		ipAddr, err := libvirt.GetDomainIPAddress(dom, image)
		if err == nil {
			log.Printf("VM is up with IP address %s\n", ipAddr)
			return ipAddr
		}

		time.Sleep(2 * time.Second)
	}

	log.Fatalf("VM is not up\n")
	return "" // Does not actually return.
}

func configureClient(c *UpCommand, dom *rawLibvirt.Domain, image *images.Image, guestIPAddr string) {
	client, err := goph.NewUnknown(image.SSHUser, guestIPAddr, goph.Password(image.SSHPassword))
	if err != nil {
		log.Fatal(err)
	}

	publicKeyPath := paths.GetDataFilePath(constants.SSHKeypairName + ".pub")
	publicKey, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatalf("Unable to read private SSH key: %s\n", err)
	}
	publicKeyStr := string(publicKey[:])

	cmds := append(image.ConfigurationCmds, []string{
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
		fmt.Sprintf("DISPLAY=:0 setxkbmap %s", host.GetHostKeyboardLayout()),
	}...)

	for _, cmd := range cmds {
		_, err := client.Run(cmd)
		if err != nil {
			log.Fatalf("Failed to run command over SSH: %s\n", err)
		}
	}
}

func ensureSSHKeypairExists() error {
	sshKeypairPath := paths.GetDataFilePath(constants.SSHKeypairName)

	if _, err := os.Stat(sshKeypairPath); err == nil {
		// SSH keypair already exists.
		return nil
	}

	cmd := exec.Command(
		paths.GetCmdPathOrExit("ssh-keygen"),
		"-f",
		sshKeypairPath,
		"-t",
		"ed25519",
		"-C",
		constants.SSHKeypairName,
		"-q",
		"-N",
		"", // Password is empty so no typing is required.
	)

	if err := cmd.Start(); err != nil {
		return err
	}

	return cmd.Wait()
}

func (c *UpCommand) Run(s *settings.Settings) {
	image := images.GetImageDetails(s.Type)

	if err := ensureSSHKeypairExists(); err != nil {
		log.Fatalf("Cannot create SSH keypair: %s\n", err)
	}

	xml := buildXML(c, image)

	conn := libvirt.Connect()
	defer conn.Close()

	dom, err := conn.DomainCreateXML(xml, 0)
	if err != nil {
		log.Fatalf("Cannot create domain: %s\n", err)
	}
	defer dom.Free()

	image.Boot(dom)
	guestIPAddr := waitBootComplete(dom, &image)
	image.StartSSH(dom)

	configureClient(c, dom, &image, guestIPAddr)

	log.Printf("%s is now ready to use\n", image.DisplayName)
}
