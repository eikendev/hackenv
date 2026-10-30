package doctor

import (
	"encoding/xml"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/sys/unix"
	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/constants"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/paths"
)

const (
	minLibvirtVersion = 9_000_000
	selinuxMount      = "/sys/fs/selinux"
	selinuxLabel      = "svirt_image_t"
	kvmDevice         = "/dev/kvm"
	keyExpiryWarning  = 90 * 24 * time.Hour
)

var hostArchToLibvirt = map[string]string{
	"amd64": "x86_64",
	"arm64": "aarch64",
}

// Checks returns all host checks in the order they should run.
func Checks() []Check {
	checks := []Check{
		{Name: "libvirt library", Run: checkLibvirtLibrary},
		{Name: "libvirt session", Run: checkLibvirtSession},
		{Name: "libvirt version", Run: checkLibvirtVersion},
		{Name: "kvm device", Run: checkKVMDevice},
		{Name: "kvm guest support", Run: checkKVMGuest},
		{Name: "passt", Run: commandCheck("passt", "passt", "install the passt package")},
		{Name: "ssh", Run: commandCheck("ssh", "ssh", "install the OpenSSH client")},
		{Name: "ssh-keygen", Run: commandCheck("ssh-keygen", "ssh-keygen", "install the OpenSSH client")},
		{Name: "viewer", Run: anyCommandCheck("viewer", []string{"virt-viewer", "remmina"}, "install virt-viewer or remmina for hackenv gui")},
		{Name: "selinux labels", Run: checkSELinuxLabels},
		{Name: "signing keys", Run: checkSigningKeys},
	}

	for _, image := range images.GetAllImages() {
		checks = append(checks, Check{Name: "image " + image.Name, Run: imageCheck(image)})
	}

	return checks
}

func connect() (*rawLibvirt.Connect, error) {
	return rawLibvirt.NewConnect(constants.ConnectURI)
}

func closeConnect(conn *rawLibvirt.Connect) {
	if conn != nil {
		_, _ = conn.Close()
	}
}

func checkLibvirtLibrary() Result {
	const name = "libvirt library"

	version, err := rawLibvirt.GetVersion()
	if err != nil {
		return fail(name, err.Error(), "install the libvirt client library (libvirt-libs or libvirt0)")
	}

	return pass(name, formatLibvirtVersion(version))
}

func checkLibvirtSession() Result {
	const name = "libvirt session"

	conn, err := connect()
	if err != nil {
		return fail(name, err.Error(), "install libvirt with the QEMU driver and make sure your user can run it")
	}
	closeConnect(conn)

	return pass(name, constants.ConnectURI)
}

func checkLibvirtVersion() Result {
	const name = "libvirt version"

	conn, err := connect()
	if err != nil {
		return skip(name, "no libvirt session")
	}
	defer closeConnect(conn)

	version, err := conn.GetLibVersion()
	if err != nil {
		return fail(name, err.Error(), "")
	}

	if version < minLibvirtVersion {
		return fail(name, formatLibvirtVersion(version), fmt.Sprintf("libvirt %s or newer is required for passt networking", formatLibvirtVersion(minLibvirtVersion)))
	}

	return pass(name, formatLibvirtVersion(version))
}

func formatLibvirtVersion(version uint32) string {
	return fmt.Sprintf("%d.%d.%d", version/1_000_000, (version/1_000)%1_000, version%1_000)
}

func checkKVMDevice() Result {
	const name = "kvm device"

	file, err := os.OpenFile(kvmDevice, os.O_RDWR, 0)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fail(name, err.Error(), "enable virtualization in the firmware and load the kvm module")
		}
		return fail(name, err.Error(), "add your user to the kvm group and log in again")
	}
	_ = file.Close()

	return pass(name, kvmDevice)
}

type capabilities struct {
	Guests []struct {
		Arch struct {
			Name     string `xml:"name,attr"`
			Emulator string `xml:"emulator"`
			Domains  []struct {
				Type string `xml:"type,attr"`
			} `xml:"domain"`
		} `xml:"arch"`
	} `xml:"guest"`
}

func checkKVMGuest() Result {
	const name = "kvm guest support"

	conn, err := connect()
	if err != nil {
		return skip(name, "no libvirt session")
	}
	defer closeConnect(conn)

	desc, err := conn.GetCapabilities()
	if err != nil {
		return fail(name, err.Error(), "")
	}

	emulator, err := findKVMEmulator(desc, hostArchToLibvirt[runtime.GOARCH])
	if err != nil {
		return fail(name, err.Error(), "install QEMU with KVM support (qemu-kvm or qemu-system)")
	}

	return pass(name, emulator)
}

func findKVMEmulator(desc, arch string) (string, error) {
	var caps capabilities
	if err := xml.Unmarshal([]byte(desc), &caps); err != nil {
		return "", fmt.Errorf("cannot parse capabilities: %w", err)
	}

	for _, guest := range caps.Guests {
		if guest.Arch.Name != arch {
			continue
		}
		for _, domain := range guest.Arch.Domains {
			if domain.Type == "kvm" {
				return guest.Arch.Emulator, nil
			}
		}
	}

	return "", fmt.Errorf("no kvm guest support for %s", arch)
}

func commandCheck(name, command, hint string) func() Result {
	return func() Result {
		path, err := exec.LookPath(command)
		if err != nil {
			return fail(name, fmt.Sprintf("%s not found in PATH", command), hint)
		}
		return pass(name, path)
	}
}

func anyCommandCheck(name string, commands []string, hint string) func() Result {
	return func() Result {
		for _, command := range commands {
			if path, err := exec.LookPath(command); err == nil {
				return pass(name, path)
			}
		}
		return fail(name, fmt.Sprintf("none of %s found in PATH", strings.Join(commands, ", ")), hint)
	}
}

func checkSELinuxLabels() Result {
	const name = "selinux labels"

	if _, err := os.Stat(selinuxMount); err != nil {
		return skip(name, "selinux not active")
	}

	sharedPath, err := paths.GetDataFilePath("shared")
	if err != nil {
		return fail(name, err.Error(), "")
	}

	targets := []string{sharedPath}
	for _, image := range images.GetAllImages() {
		if path, err := image.GetLatestPath(); err == nil {
			targets = append(targets, path)
		}
	}

	var unlabeled []string
	for _, target := range targets {
		if !hasSELinuxType(target, selinuxLabel) {
			unlabeled = append(unlabeled, target)
		}
	}

	if len(unlabeled) > 0 {
		return fail(name, "missing "+selinuxLabel+" on "+strings.Join(unlabeled, ", "), "run hackenv fix apply-labels")
	}

	return pass(name, selinuxLabel)
}

func hasSELinuxType(path, wantType string) bool {
	buf := make([]byte, 256)
	n, err := unix.Getxattr(path, "security.selinux", buf)
	if err != nil {
		return false
	}

	label := strings.TrimRight(string(buf[:n]), "\x00")
	parts := strings.Split(label, ":")
	if len(parts) < 3 {
		return false
	}

	return parts[2] == wantType
}

func checkSigningKeys() Result {
	const name = "signing keys"

	return checkKeyExpiry(name, images.SigningKeys(), time.Now())
}

func checkKeyExpiry(name string, keys []images.SigningKey, now time.Time) Result {
	var expiring []string

	for _, key := range keys {
		expiry, err := key.Expiry()
		if err != nil {
			return fail(name, fmt.Sprintf("%s: %v", key.Image, err), "")
		}
		if expiry.IsZero() {
			continue
		}
		if now.After(expiry) {
			return fail(name, fmt.Sprintf("%s key expired on %s", key.Image, expiry.Format(time.DateOnly)), "update hackenv or pass --signing-key to hackenv get")
		}
		if expiry.Sub(now) < keyExpiryWarning {
			expiring = append(expiring, fmt.Sprintf("%s expires %s", key.Image, expiry.Format(time.DateOnly)))
		}
	}

	if len(expiring) > 0 {
		return fail(name, strings.Join(expiring, ", "), "update hackenv once a release with the new key is available")
	}

	return pass(name, "valid")
}

func imageCheck(image images.Image) func() Result {
	return func() Result {
		name := "image " + image.Name

		path, err := image.GetLatestPath()
		if err != nil {
			return fail(name, "not downloaded", fmt.Sprintf("run hackenv --type=%s get", image.Name))
		}

		return pass(name, path)
	}
}
