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

The **preferred way** is to download the binary from the [latest release](https://github.com/eikendev/hackenv/releases).

Alternatively, install the [required dependencies](#dependencies) and build it yourself:
```bash
go install github.com/eikendev/hackenv/cmd/...@latest
```

Ensure that the libvirt daemon is running, and a socket created. On Fedora, this can be done as follows:
```
sudo systemctl start libvirtd
sudo systemctl start virtqemud.socket
```

## 🤘&nbsp;Features

- Support for Kali Linux and Parrot Security.
- Create short-lived virtual machines within seconds.
- Download the latest official live image conveniently.
- **Simple and intuitive** command line interface.
- Configure instant SSH access based on **public-key authentication**.
- Set up a **shared directory** between host and guest.
- Set the same **keyboard layout** in the guest as on the host.

## 📄&nbsp;Usage

First, make sure you have the [required dependencies](#dependencies) installed.

Then, download an image using `hackenv get`.
This will download a live image from the official mirrors.
The download can take a while, so sit back and enjoy some tea.

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

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
You can run `hackenv fix apply-labels` (or `./bin/hackenv_applylabels`) if you are on Fedora or similar.
The script will also relabel the downloaded images.
If you add new files to the shared directory from outside, do not forget to label them as well.

## ⚙&nbsp;Configuration

The tool currently does not support configuration via files.
However, some options can be set using environment variables.
Check out the help (`--help`) of a command to see what options support this.

For instance, to operate with Parrot Security by default, you can set `$HACKENV_TYPE=parrot`.
If you work with both operating systems, then I recommend using shell aliases:
```bash
alias kali='hackenv --type=kali'
alias parrot='hackenv --type=parrot'
```

## 🥙&nbsp;Dependencies

- [libvirt](https://libvirt.org/) 9.0 or newer with [QEMU](https://www.qemu.org/) and KVM
- [passt](https://passt.top/)
- [OpenSSH](https://www.openssh.com/) (ssh and ssh-keygen)
- setxkbmap
- [virt-viewer](https://virt-manager.org/) or [Remmina](https://remmina.org/)

The binary loads the libvirt client library at runtime, which every libvirt installation provides (`libvirt-libs` on Fedora, `libvirt0` on Debian and Ubuntu).
Building from source needs a C compiler but no libvirt development files.

## 💡&nbsp;Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
