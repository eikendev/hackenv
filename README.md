<div align="center">
	<h1>hackenv</h1>
	<h4 align="center">
		Access your <a href="https://www.kali.org/">Kali Linux</a> and <a href="https://parrotsec.org/">Parrot Security</a> instances with ease.
	</h4>
	<p>hackenv lets you comfortably manage your hacking environment from the terminal.</p>
</div>

<p align="center">
	<a href="https://github.com/eikendev/hackenv/actions"><img alt="Build status" src="https://img.shields.io/github/actions/workflow/status/eikendev/hackenv/release.yaml?branch=main"/></a>&nbsp;
	<a href="https://github.com/eikendev/hackenv/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/eikendev/hackenv"/></a>&nbsp;
</p>

## 🤘&nbsp;Features

- Support for Kali Linux and Parrot Security.
- Create short-lived virtual machines within seconds.
- Download the latest official live image conveniently, verified against the vendor's signature.
- **Simple and intuitive** command line interface.
- Configure instant SSH access based on **public-key authentication**.
- Set up a **shared directory** between host and guest.
- No root privileges, no system daemon, no network bridge: everything runs in your user session.

## 🚀&nbsp;Installation

The **preferred way** is to download the binary for amd64 or arm64 from the [latest release](https://github.com/eikendev/hackenv/releases).
Alternatively, [build it yourself](#building-from-source).

hackenv needs [libvirt](https://libvirt.org/) 9.0 or newer, [QEMU](https://www.qemu.org/) with KVM, [passt](https://passt.top/), [OpenSSH](https://www.openssh.com/), and [virt-viewer](https://virt-manager.org/) or [Remmina](https://remmina.org/) for the GUI.
Install them for your distribution:

| Distribution | Install dependencies |
|---|---|
| Fedora 40+ | `sudo dnf install libvirt-daemon-kvm qemu-kvm passt virt-viewer openssh-clients` |
| Debian 12+, LMDE 6+ | `sudo apt install libvirt-daemon libvirt-clients qemu-system-x86 passt virt-viewer openssh-client` |
| Ubuntu 24.04+, Linux Mint 22+ | `sudo apt install libvirt-daemon libvirt-clients qemu-system-x86 passt virt-viewer openssh-client` |
| Arch Linux and derivatives | `sudo pacman -S libvirt qemu-desktop passt virt-viewer openssh` |
| RHEL 9.2+, AlmaLinux, Rocky Linux, CentOS Stream 9+ | `sudo dnf install libvirt qemu-kvm passt virt-viewer openssh-clients` |

On Debian, Ubuntu and their derivatives, add yourself to the `kvm` group and log in again so that your user can open `/dev/kvm`:
```bash
sudo usermod -aG kvm "$USER"
```

Finally, run `hackenv doctor`.
It checks every prerequisite and tells you what is missing.

> [!NOTE]
> openSUSE Tumbleweed and Leap 15.6, Pop!_OS 24.04 and Gentoo ship recent enough versions of libvirt and passt according to their package lists, but have not been tested.
> Ubuntu 22.04, Linux Mint 21, Debian 11 and RHEL 8 ship a libvirt older than 9.0 and are not supported.
> Alpine Linux (musl) and NixOS (no FHS library paths) cannot run the release binary; build from source there.

## ⚡&nbsp;Quick Start

```bash
hackenv get   # download and verify the latest Kali Linux image
hackenv up    # boot the VM and set up SSH and the shared directory
hackenv ssh   # open a shell in the VM
hackenv gui   # or open its desktop
hackenv down  # shut the VM down and discard its state
```

Add `--type=parrot` to any command to work with Parrot Security instead.

## 📄&nbsp;Usage

### Choosing a Distribution

hackenv operates with Kali Linux by default.
Pass `--type=parrot` to work with Parrot Security instead, or set it permanently through `$HACKENV_TYPE` or the [configuration file](#configuration).
Both distributions can be used side by side; each has its own image, VM and SSH port.

### Downloading an Image

`hackenv get` downloads the latest live image into `~/.local/share/hackenv`.
Parrot Security is fetched from the official mirrors over HTTPS.
hackenv picks Parrot's *tuned* image, which is built for x86-64-v3 CPUs (Intel Haswell, AMD Ryzen and newer); on an older CPU, download the regular `Parrot-security-<version>_amd64.iso` yourself and import it with `--file`.
Kali Linux publishes its live images only as torrents, so hackenv downloads them with a built-in BitTorrent client: it only leeches, uses a random port, and exits as soon as the image is complete.
The download can take a while, so sit back and enjoy some tea.
An interrupted torrent download resumes when you run the command again.

Before anything is downloaded, hackenv verifies the vendor's OpenPGP signature on the checksum file using the public keys embedded in the binary, and then checks the image against that file.

| Vendor | Signing key fingerprint |
|---|---|
| Kali Linux | `827C 8569 F251 8CC6 77FE CA1A ED65 462E C8D5 E4C5` |
| Parrot Security | `B711 8223 4655 2E4D 92DA 02DF 7A82 86AF 0E81 EE4A` |

| Flag | Purpose |
|---|---|
| `--update` | Download the latest release when you already have an older image. |
| `--force` | Download even if the latest image is already present. |
| `--yes` | Answer yes to all questions. |
| `--file PATH` | Import an image you already have instead of downloading one, for example an older release or the regular Parrot Security image. It is copied as is and **not verified**, so make sure it is authentic. Keep the vendor's filename, hackenv reads the release from it. |
| `--signing-key FILE` | Verify with this armored public key instead of the embedded one, for when a vendor rotates its key before hackenv ships the new one. |
| `--skip-signature` | Bypass the signature check entirely. Meant as a last resort. |

> [!WARNING]
> If you run SELinux, label the new image before booting it: `hackenv fix apply-labels`.
> See [Troubleshooting](#troubleshooting) for details.

### Starting the VM

`hackenv up` boots the virtual machine.
Once this command is finished, the VM is running and fully configured.
The VM is short-lived (volatile), meaning any data and configuration stored outside the shared directory is removed once the machine shuts down.
This is by design and admittedly opinionated.

| Flag | Default | Purpose |
|---|---|---|
| `--memory` | `4G` | RAM of the VM, for example `8G` or `4096M`. Capped at the memory of your host. |
| `--cores` | `2` | Number of virtual CPU cores. |
| `--ssh-port` | `2222` (Kali), `2223` (Parrot) | Host port forwarded to the SSH server of the VM. Pick another one if it is taken. |

The VM uses user-mode networking via [passt](https://passt.top/), so no bridge or privileges are needed.
Its SSH server is reachable at `127.0.0.1` on the port above.
`hackenv up` logs which host features it detected (machine type, CPU mode, video model) and picks sensible fallbacks when a modern one is missing.

### Connecting to the VM

- `hackenv ssh` opens an SSH session using the key hackenv generated for you.
- `hackenv gui` opens the desktop of the VM in virt-viewer (or Remmina with `--viewer=remmina`); add `--fullscreen` if you like.
- `hackenv status` prints whether the VM is running.
- `hackenv down` shuts the VM down and discards its state.

### Sharing Files

hackenv automatically sets up a shared directory between the host and the virtual machine.
On the host side the directory is `~/.local/share/hackenv/shared`, while on the guest side it is located at `/shared`.
It is shared through 9p in mapped mode, so everything the guest writes belongs to your host user.

### Troubleshooting

`hackenv doctor` checks every prerequisite: the libvirt library and session, KVM, passt, OpenSSH, the viewer, SELinux labels, the expiry of the embedded signing keys, and whether an image is present.
Run it first whenever something does not work.

If SELinux denies access to the image or the shared directory, you have to adjust the context of those files.
`hackenv fix apply-labels` (or `./bin/hackenv_applylabels`) relabels the downloaded images and the shared directory, and `hackenv doctor` reports what still lacks the label.
If you add new files to the shared directory from outside, do not forget to label them as well.
On hosts without SELinux the command does nothing.

## ⚙&nbsp;Configuration

Every flag can also be set through an environment variable or a YAML config file.
Check out the help (`--help`) of a command to see the variable names.
Command line flags win over environment variables, which win over the config file.

The config file lives at `~/.config/hackenv/config.yaml` (or wherever `$XDG_CONFIG_HOME` points), and `--config FILE` selects another one.
Global flags are top-level keys, flags of a command are nested under the command name:
```yaml
type: parrot
up:
  memory: 8G
  cores: 4
  ssh-port: 2300
```

If you work with both operating systems, then I recommend using shell aliases:
```bash
alias kali='hackenv --type=kali'
alias parrot='hackenv --type=parrot'
```

## 🔨&nbsp;Building from Source

```bash
go install -tags libvirt_dlopen,nosqlite,disable_libutp github.com/eikendev/hackenv/cmd/...@latest
```

Building needs a C compiler and the libvirt development headers (`libvirt-devel` on Fedora, `libvirt-dev` on Debian and Ubuntu).
The resulting binary loads the libvirt client library at runtime, which every libvirt installation provides (`libvirt-libs` on Fedora, `libvirt0` on Debian and Ubuntu), so one build runs on any libvirt 9.0 or newer.

## 💡&nbsp;Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
