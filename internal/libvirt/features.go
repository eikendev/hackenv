package libvirt

import (
	"encoding/xml"
	"fmt"
	"log/slog"
	"runtime"
	"slices"

	"libvirt.org/go/libvirt"
)

const minVirtiofsSessionVersion = 10_000_000

var hostArchToLibvirt = map[string]string{
	"amd64": "x86_64",
	"arm64": "aarch64",
}

var preferredMachines = map[string]string{
	"x86_64":  "q35",
	"aarch64": "virt",
}

// HostFeatures lists what the local hypervisor supports for the VM definition.
type HostFeatures struct {
	Machine    string
	CPUMode    string
	VideoModel string
	Virtiofs   bool
}

// ConservativeFeatures returns features every KVM host supports.
func ConservativeFeatures() HostFeatures {
	return HostFeatures{VideoModel: "qxl"}
}

type capabilitiesXML struct {
	Guests []struct {
		Arch struct {
			Name     string `xml:"name,attr"`
			Machines []struct {
				Name      string `xml:"name,attr"`
				Canonical string `xml:"canonical,attr"`
				Value     string `xml:",chardata"`
			} `xml:"machine"`
		} `xml:"arch"`
	} `xml:"guest"`
}

type domainCapabilitiesXML struct {
	CPU struct {
		Modes []struct {
			Name      string `xml:"name,attr"`
			Supported string `xml:"supported,attr"`
		} `xml:"mode"`
	} `xml:"cpu"`
	Devices struct {
		Video struct {
			Enums []enumXML `xml:"enum"`
		} `xml:"video"`
		Filesystem struct {
			Enums []enumXML `xml:"enum"`
		} `xml:"filesystem"`
	} `xml:"devices"`
}

type enumXML struct {
	Name   string   `xml:"name,attr"`
	Values []string `xml:"value"`
}

func enumValues(enums []enumXML, name string) []string {
	for _, e := range enums {
		if e.Name == name {
			return e.Values
		}
	}
	return nil
}

// DetectHostFeatures inspects the hypervisor and picks the best supported options.
// Any detection failure degrades to conservative defaults.
func DetectHostFeatures(conn *libvirt.Connect) HostFeatures {
	features := ConservativeFeatures()
	arch := hostArchToLibvirt[runtime.GOARCH]

	caps, err := conn.GetCapabilities()
	if err != nil {
		slog.Warn("failed to read host capabilities", "err", err)
		return features
	}
	features.Machine = pickMachine(caps, arch)

	domCaps, err := conn.GetDomainCapabilities("", arch, features.Machine, "kvm", 0)
	if err != nil {
		slog.Warn("failed to read domain capabilities", "err", err)
		return features
	}

	version, err := conn.GetLibVersion()
	if err != nil {
		slog.Warn("failed to read libvirt version", "err", err)
		version = 0
	}

	return applyDomainCapabilities(features, domCaps, version)
}

func pickMachine(caps, arch string) string {
	var parsed capabilitiesXML
	if err := xml.Unmarshal([]byte(caps), &parsed); err != nil {
		slog.Warn("failed to parse host capabilities", "err", err)
		return ""
	}

	preferred := preferredMachines[arch]
	for _, guest := range parsed.Guests {
		if guest.Arch.Name != arch {
			continue
		}
		for _, machine := range guest.Arch.Machines {
			if machine.Value == preferred {
				return preferred
			}
		}
	}

	return ""
}

func applyDomainCapabilities(features HostFeatures, domCaps string, libvirtVersion uint32) HostFeatures {
	var parsed domainCapabilitiesXML
	if err := xml.Unmarshal([]byte(domCaps), &parsed); err != nil {
		slog.Warn("failed to parse domain capabilities", "err", err)
		return features
	}

	for _, mode := range parsed.CPU.Modes {
		if mode.Supported != "yes" {
			continue
		}
		if mode.Name == "host-passthrough" {
			features.CPUMode = "host-passthrough"
			break
		}
		if mode.Name == "host-model" {
			features.CPUMode = "host-model"
		}
	}

	if slices.Contains(enumValues(parsed.Devices.Video.Enums, "modelType"), "virtio") {
		features.VideoModel = "virtio"
	}

	if libvirtVersion >= minVirtiofsSessionVersion && slices.Contains(enumValues(parsed.Devices.Filesystem.Enums, "driverType"), "virtiofs") {
		features.Virtiofs = true
	}

	return features
}

// String summarises the features for logging.
func (f HostFeatures) String() string {
	return fmt.Sprintf("machine=%q cpu=%q video=%q virtiofs=%t", f.Machine, f.CPUMode, f.VideoModel, f.Virtiofs)
}
