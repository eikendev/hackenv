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

## 🚀&nbsp;Installation

The **preferred way** is to download the binary for amd64 or arm64 from the [latest release](https://github.com/eikendev/hackenv/releases).

Alternatively, build it yourself:
```bash
go install -tags libvirt_dlopen github.com/eikendev/hackenv/cmd/...@latest
```

Then install the [dependencies](#dependencies) for your distribution and run `hackenv doctor`.
It checks every prerequisite and tells you what is missing.

### Supported Platforms

hackenv needs libvirt 9.0 or newer, QEMU with KVM, and passt.
It talks to libvirt over the unprivileged `qemu:///session` connection, so no system daemon has to be configured and no root privileges are needed.

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

> [!NOTE]
> openSUSE Tumbleweed and Leap 15.6, Pop!_OS 24.04 and Gentoo ship recent enough versions of libvirt and passt according to their package lists, but have not been tested.

Ubuntu 22.04, Linux Mint 21, Debian 11 and RHEL 8 ship a libvirt older than 9.0 and are not supported.
Alpine Linux (musl) and NixOS (no FHS library paths) cannot run the release binary; build from source there.

## 🤘&nbsp;Features

- Support for Kali Linux and Parrot Security.
- Create short-lived virtual machines within seconds.
- Download the latest official live image conveniently.
- **Simple and intuitive** command line interface.
- Configure instant SSH access based on **public-key authentication**.
- Set up a **shared directory** between host and guest.

## 📄&nbsp;Usage

First, make sure you have the [required dependencies](#dependencies) installed.

Then, download an image using `hackenv get`.
Parrot Security is fetched from the official mirrors over HTTPS.
Kali Linux publishes its live images only as torrents, so hackenv downloads them with a built-in BitTorrent client: it only leeches, uses a random port, and exits as soon as the image is complete.
The download can take a while, so sit back and enjoy some tea.

Before anything is downloaded, hackenv verifies the vendor's OpenPGP signature on the checksum file using the public keys embedded in the binary (Kali `827C 8569 F251 8CC6 77FE CA1A ED65 462E C8D5 E4C5`, Parrot `B711 8223 4655 2E4D 92DA 02DF 7A82 86AF 0E81 EE4A`), and then checks the image against that file.
If you already have the image, for example from the official torrent, import it with `hackenv get --file PATH`; it is checked against the same signed checksum list before it is copied into place.
If a vendor rotates its key before hackenv ships the new one, pass the fresh key with `hackenv get --signing-key FILE`.
`hackenv get --skip-signature` bypasses the signature check entirely and is meant as a last resort.

| :warning:  **If you run SELinux, you must label the new image yourself. Check [the section on file sharing](#file-sharing) for more information.**   |
|-----------------------------------------------------------------|

Next, run `hackenv up` to boot the virtual machine.
Once this command is finished, the VM is running and fully configured.
The VM will be short-lived (volatile), meaning any data and configuration stored outside the shared directory will be removed once the machine shuts down.
This is by design and admittedly opinionated.

You can now start an SSH session with `hackenv ssh` or spin up a GUI with `hackenv gui`.
The VM uses user-mode networking via [passt](https://passt.top/), so no bridge or privileges are needed.
Its SSH server is forwarded to `127.0.0.1:2222` for Kali Linux and `127.0.0.1:2223` for Parrot Security.
If a port is taken on your host, pick another one with `hackenv up --ssh-port` or `$HACKENV_SSH_PORT`.

Note that by default, hackenv will operate with Kali Linux, and respectively download its image.
If you want to operate with Parrot Security instead, specify `hackenv --type=parrot`, or check out [the configuration](#configuration).

### File Sharing

hackenv will automatically set up a shared directory between the host and the virtual machine.
On the host side the directory is `~/.local/share/hackenv/shared`, while on the guest side it is located at `/shared`.
It is shared through 9p in mapped mode, so everything the guest writes belongs to your host user; `hackenv up` logs which other host features it detected.

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
You can run `hackenv fix apply-labels` (or `./bin/hackenv_applylabels`) if you are on Fedora or similar.
The script will also relabel the downloaded images, and `hackenv doctor` reports images and directories that still lack the label.
If you add new files to the shared directory from outside, do not forget to label them as well.
On hosts without SELinux the script does nothing.

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

For instance, to operate with Parrot Security by default, you can also set `$HACKENV_TYPE=parrot`.
If you work with both operating systems, then I recommend using shell aliases:
```bash
alias kali='hackenv --type=kali'
alias parrot='hackenv --type=parrot'
```

## 🥙&nbsp;Dependencies

- [libvirt](https://libvirt.org/) 9.0 or newer with [QEMU](https://www.qemu.org/) and KVM
- [passt](https://passt.top/)
- [OpenSSH](https://www.openssh.com/) (ssh and ssh-keygen)
- [virt-viewer](https://virt-manager.org/) or [Remmina](https://remmina.org/) for `hackenv gui`

Run `hackenv doctor` to verify that everything is in place.

The binary loads the libvirt client library at runtime, which every libvirt installation provides (`libvirt-libs` on Fedora, `libvirt0` on Debian and Ubuntu), so one release runs on any libvirt 9.0 or newer.
Building from source needs a C compiler and the libvirt development headers (`libvirt-devel` on Fedora, `libvirt-dev` on Debian and Ubuntu).

## 💡&nbsp;Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
