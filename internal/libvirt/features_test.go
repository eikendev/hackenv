package libvirt

import "testing"

const sampleCaps = `<capabilities>
  <guest>
    <os_type>hvm</os_type>
    <arch name='x86_64'>
      <machine maxCpus='255'>pc-i440fx-9.0</machine>
      <machine canonical='pc-i440fx-9.0' maxCpus='255'>pc</machine>
      <machine maxCpus='1024'>pc-q35-9.0</machine>
      <machine canonical='pc-q35-9.0' maxCpus='1024'>q35</machine>
    </arch>
  </guest>
</capabilities>`

const sampleDomCaps = `<domainCapabilities>
  <cpu>
    <mode name='host-passthrough' supported='yes'/>
    <mode name='maximum' supported='yes'/>
    <mode name='host-model' supported='yes'/>
  </cpu>
  <devices>
    <video supported='yes'>
      <enum name='modelType'>
        <value>vga</value>
        <value>qxl</value>
        <value>virtio</value>
      </enum>
    </video>
    <filesystem supported='yes'>
      <enum name='driverType'>
        <value>path</value>
        <value>handle</value>
        <value>virtiofs</value>
      </enum>
    </filesystem>
  </devices>
</domainCapabilities>`

const oldDomCaps = `<domainCapabilities>
  <cpu>
    <mode name='host-passthrough' supported='no'/>
    <mode name='host-model' supported='yes'/>
  </cpu>
  <devices>
    <video supported='yes'>
      <enum name='modelType'>
        <value>vga</value>
        <value>qxl</value>
      </enum>
    </video>
    <filesystem supported='yes'>
      <enum name='driverType'>
        <value>path</value>
      </enum>
    </filesystem>
  </devices>
</domainCapabilities>`

func TestPickMachine(t *testing.T) {
	if got := pickMachine(sampleCaps, "x86_64"); got != "q35" {
		t.Fatalf("pickMachine() = %q, want q35", got)
	}
	if got := pickMachine(sampleCaps, "aarch64"); got != "" {
		t.Fatalf("pickMachine() for missing arch = %q, want empty", got)
	}
	if got := pickMachine("<capabilities>", "x86_64"); got != "" {
		t.Fatalf("pickMachine() for malformed xml = %q, want empty", got)
	}
}

func TestApplyDomainCapabilities(t *testing.T) {
	modern := applyDomainCapabilities(ConservativeFeatures(), sampleDomCaps, 12_000_000)
	if modern.CPUMode != "host-passthrough" || modern.VideoModel != "virtio" || !modern.Virtiofs {
		t.Fatalf("modern host = %s", modern)
	}

	oldLibvirt := applyDomainCapabilities(ConservativeFeatures(), sampleDomCaps, 9_000_000)
	if oldLibvirt.Virtiofs {
		t.Fatalf("virtiofs enabled on libvirt 9: %s", oldLibvirt)
	}

	old := applyDomainCapabilities(ConservativeFeatures(), oldDomCaps, 12_000_000)
	if old.CPUMode != "host-model" || old.VideoModel != "qxl" || old.Virtiofs {
		t.Fatalf("old host = %s", old)
	}

	broken := applyDomainCapabilities(ConservativeFeatures(), "<domainCapabilities>", 12_000_000)
	if broken != ConservativeFeatures() {
		t.Fatalf("malformed caps changed features: %s", broken)
	}
}
