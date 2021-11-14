<div align="center">
	<h1>hackenv</h1>
	<h4 align="center">
		Access your <a href="https://www.kali.org/">Kali Linux</a> and <a href="https://parrotsec.org/">Parrot Security</a> instances with ease.
	</h4>
	<p>hackenv lets you comfortably manage your hacking environment from the terminal.</p>
</div>

<p align="center">
	<a href="https://github.com/eikendev/hackenv/actions"><img alt="Build status" src="https://img.shields.io/github/workflow/status/eikendev/hackenv/Main"/></a>&nbsp;
	<a href="https://github.com/eikendev/hackenv/blob/master/LICENSE"><img alt="License" src="https://img.shields.io/github/license/eikendev/hackenv"/></a>&nbsp;
</p>

## 🚀&nbsp;Installation

The **preferred way** is to download the binary from the [latest release](https://github.com/eikendev/hackenv/releases).

Alternatively, install the [required dependencies](#dependencies) and build it yourself:
```bash
go install github.com/eikendev/hackenv/cmd/...@latest
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
Also, you will need a bridge interface [as described below](#creating-a-bridge-interface).
This can be as simple as running `./bin/hackenv_createbridge`.

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

Note that by default, hackenv will operate with Kali Linux, and respectively download its image.
If you want to operate with Parrot Security instead, specify `hackenv --type=parrot`, or check out [the configuration](#configuration).

### File Sharing

hackenv will automatically set up a shared directory between the host and the virtual machine.
On the host side the directory is `~/.local/share/hackenv/shared`, while on the guest side it is located at `/shared`.

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
You can run `./bin/hackenv_fixlabels` if you are on Fedora or similar.
The script will also relabel the downloaded images.
If you add new files to the shared directory from outside, do not forget to label them as well.

### Creating a Bridge Interface

hackenv uses a bridge so that you can reach the guest from the host for SSH, while the guest can access the Internet.
You can create this bridge by running `./bin/hackenv_createbridge`.
Note that this script **will request privileges** so it can create an interface.

Of course, please adapt the script to your specific needs.
The interface is expected to have the name `virbr0` by default, but this can be changed using command line flags.

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

- [libvirt](https://libvirt.org/) (virsh)
- [OpenSSH](https://www.openssh.com/) (ssh and ssh-keygen)
- setxkbmap
- [virt-viewer](https://virt-manager.org/) or [Remmina](https://remmina.org/)

To build the binary yourself, you also need the development files of libvirt, usually called `libvirt-dev` or `libvirt-devel`.

## 💡&nbsp;Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
