// Package domain renders the libvirt domain definition for a VM.
package domain

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"text/template"
)

// SharedTag is the mount tag the guest uses to mount the shared directory.
const SharedTag = "shared"

// Spec describes the VM to define.
type Spec struct {
	Name        string
	MemoryKiB   int64
	Cores       int
	ImagePath   string
	SharedPath  string
	Machine     string
	CPUMode     string
	VideoModel  string
	MacAddress  string
	ForwardHost string
	SSHPort     int
}

var domainTemplate = template.Must(template.New("domain").Funcs(template.FuncMap{"esc": escape}).Parse(`<domain type='kvm'>
  <name>{{esc .Name}}</name>
  <memory unit='KiB'>{{.MemoryKiB}}</memory>
  <vcpu placement='static'>{{.Cores}}</vcpu>
  <os>
    <type{{if .Machine}} machine='{{esc .Machine}}'{{end}}>hvm</type>
    <boot dev='cdrom'/>
  </os>
  <features>
    <acpi/>
    <apic/>
  </features>
{{- if .CPUMode}}
  <cpu mode='{{esc .CPUMode}}' check='none'{{if eq .CPUMode "host-passthrough"}} migratable='off'{{end}}/>
{{- end}}
  <devices>
    <disk type='file' device='cdrom'>
      <driver name='qemu' type='raw'/>
      <source file='{{esc .ImagePath}}'/>
      <target dev='sda' bus='sata'/>
      <readonly/>
    </disk>
    <filesystem type='mount' accessmode='mapped'>
      <source dir='{{esc .SharedPath}}'/>
      <target dir='shared'/>
    </filesystem>
    <interface type='user'>
      <mac address='{{esc .MacAddress}}'/>
      <model type='virtio'/>
      <backend type='passt'/>
      <portForward proto='tcp' address='{{esc .ForwardHost}}'>
        <range start='{{.SSHPort}}' to='22'/>
      </portForward>
    </interface>
    <console type='pty'>
      <target type='serial'/>
    </console>
    <graphics type='spice' port='-1' autoport='yes'>
      <listen type='address' address='127.0.0.1'/>
      <image compression='off'/>
    </graphics>
    <channel type='spicevmc'>
      <target type='virtio' name='com.redhat.spice.0'/>
    </channel>
    <video>
      <model type='{{esc .VideoModel}}'/>
    </video>
    <sound model='ich9'/>
    <input type='tablet' bus='virtio'/>
    <input type='keyboard' bus='ps2'/>
    <rng model='virtio'>
      <backend model='random'>/dev/urandom</backend>
    </rng>
  </devices>
</domain>
`))

func escape(s string) string {
	var buf bytes.Buffer
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		return ""
	}
	return buf.String()
}

// Render produces the domain XML for the spec.
func Render(spec Spec) (string, error) {
	if spec.VideoModel == "" {
		spec.VideoModel = "qxl"
	}

	var out bytes.Buffer
	if err := domainTemplate.Execute(&out, spec); err != nil {
		return "", fmt.Errorf("cannot render domain definition: %w", err)
	}

	return out.String(), nil
}

// MountCommand returns the guest command that mounts the shared directory at target.
func MountCommand(target string) string {
	return fmt.Sprintf("sudo mount -t 9p -o trans=virtio,version=9p2000.L %s %s", SharedTag, target)
}
