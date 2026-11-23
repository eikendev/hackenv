package domain

import (
	"encoding/xml"
	"strings"
	"testing"
)

func baseSpec() Spec {
	return Spec{
		Name:        "kali",
		MemoryKiB:   4194304,
		Cores:       2,
		ImagePath:   "/data/kali-2026.3.iso",
		SharedPath:  "/data/shared",
		MacAddress:  "52:54:00:08:f9:e8",
		ForwardHost: "127.0.0.1",
		SSHPort:     2222,
	}
}

func render(t *testing.T, spec Spec) string {
	t.Helper()

	out, err := Render(spec)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	var doc struct{}
	if err := xml.Unmarshal([]byte(out), &doc); err != nil {
		t.Fatalf("Render() produced malformed XML: %v\n%s", err, out)
	}

	return out
}

func TestRenderDefaults(t *testing.T) {
	out := render(t, baseSpec())

	for _, want := range []string{
		"<type>hvm</type>",
		"accessmode='mapped'",
		"<model type='qxl'/>",
		"<mouse mode='server'/>",
		"<range start='2222' to='22'/>",
		"<target dir='shared'/>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}

	if strings.Contains(out, "<cpu") {
		t.Errorf("unexpected cpu element in\n%s", out)
	}
}

func TestRenderModernFeatures(t *testing.T) {
	spec := baseSpec()
	spec.Machine = "q35"
	spec.CPUMode = "host-passthrough"
	spec.VideoModel = "virtio"

	out := render(t, spec)

	for _, want := range []string{
		"<type machine='q35'>hvm</type>",
		"<cpu mode='host-passthrough' check='none' migratable='off'/>",
		"<model type='virtio'/>",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in\n%s", want, out)
		}
	}
}

func TestRenderEscapesPaths(t *testing.T) {
	spec := baseSpec()
	spec.SharedPath = "/data/it's <shared> & more"

	out := render(t, spec)

	if !strings.Contains(out, "<source dir='/data/it&#39;s &lt;shared&gt; &amp; more'/>") {
		t.Fatalf("path not escaped in\n%s", out)
	}
}

func TestMountCommand(t *testing.T) {
	got := MountCommand("/shared")
	if !strings.Contains(got, "-t 9p") || !strings.HasSuffix(got, "shared /shared") {
		t.Fatalf("MountCommand() = %q", got)
	}
}

func TestAutoResizeCommands(t *testing.T) {
	got := AutoResizeCommands("kali")
	if len(got) != 2 {
		t.Fatalf("AutoResizeCommands() returned %d commands, want 2", len(got))
	}

	want := `printf '%s\n' 'ACTION=="change", SUBSYSTEM=="drm", ` +
		`RUN+="/usr/sbin/runuser -u kali -- /bin/sh -c '\''export DISPLAY=:0; xrandr >/dev/null; xrandr --output Virtual-1 --auto'\''"' ` +
		`| sudo tee /etc/udev/rules.d/50-hackenv-resize.rules >/dev/null`
	if got[0] != want {
		t.Fatalf("AutoResizeCommands()[0] = %q, want %q", got[0], want)
	}

	if got[1] != "sudo udevadm control --reload-rules" {
		t.Fatalf("AutoResizeCommands()[1] = %q", got[1])
	}
}
