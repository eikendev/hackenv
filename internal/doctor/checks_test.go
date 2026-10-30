package doctor

import (
	"testing"
	"time"

	"github.com/eikendev/hackenv/internal/images"
)

const sampleCapabilities = `<capabilities>
  <host><cpu><arch>x86_64</arch></cpu></host>
  <guest>
    <os_type>hvm</os_type>
    <arch name='i686'>
      <emulator>/usr/bin/qemu-system-i386</emulator>
      <domain type='qemu'/>
    </arch>
  </guest>
  <guest>
    <os_type>hvm</os_type>
    <arch name='x86_64'>
      <emulator>/usr/bin/qemu-system-x86_64</emulator>
      <domain type='qemu'/>
      <domain type='kvm'/>
    </arch>
  </guest>
</capabilities>`

func TestFindKVMEmulator(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		arch    string
		want    string
		wantErr bool
	}{
		{name: "kvm available", desc: sampleCapabilities, arch: "x86_64", want: "/usr/bin/qemu-system-x86_64"},
		{name: "arch without kvm", desc: sampleCapabilities, arch: "i686", wantErr: true},
		{name: "unknown arch", desc: sampleCapabilities, arch: "aarch64", wantErr: true},
		{name: "malformed xml", desc: "<capabilities>", arch: "x86_64", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findKVMEmulator(tt.desc, tt.arch)
			if (err != nil) != tt.wantErr {
				t.Fatalf("findKVMEmulator() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("findKVMEmulator() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestFormatLibvirtVersion(t *testing.T) {
	tests := map[uint32]string{
		9_000_000:  "9.0.0",
		10_000_000: "10.0.0",
		11_002_003: "11.2.3",
	}

	for version, want := range tests {
		if got := formatLibvirtVersion(version); got != want {
			t.Errorf("formatLibvirtVersion(%d) = %q, want %q", version, got, want)
		}
	}
}

func TestCheckKeyExpiry(t *testing.T) {
	keys := images.SigningKeys()
	expiry, err := keys[0].Expiry()
	if err != nil {
		t.Fatalf("expiry: %v", err)
	}

	tests := []struct {
		name string
		now  time.Time
		want Status
	}{
		{name: "long before expiry", now: expiry.Add(-365 * 24 * time.Hour), want: Pass},
		{name: "shortly before expiry", now: expiry.Add(-24 * time.Hour), want: Fail},
		{name: "after expiry", now: expiry.Add(24 * time.Hour), want: Fail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkKeyExpiry("signing keys", keys[:1], tt.now); got.Status != tt.want {
				t.Fatalf("checkKeyExpiry() = %s (%s), want %s", got.Status, got.Detail, tt.want)
			}
		})
	}
}

func TestHasFailures(t *testing.T) {
	if HasFailures([]Result{{Status: Pass}, {Status: Skip}}) {
		t.Fatal("HasFailures() = true without failures")
	}
	if !HasFailures([]Result{{Status: Pass}, {Status: Fail}}) {
		t.Fatal("HasFailures() = false with a failure")
	}
}

func TestStatusString(t *testing.T) {
	if Pass.String() != "PASS" || Fail.String() != "FAIL" || Skip.String() != "SKIP" {
		t.Fatal("unexpected status strings")
	}
}
