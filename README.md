<div align="center">
	<h1>hackenv</h1>
	<h4 align="center">
		Access your <a href="https://www.kali.org/">Kali Linux</a> and <a href="https://www.parrotsec.org/">Parrot Security</a> instances with ease.
	</h4>
	<p>hackenv lets you comfortably manage your security environment from the command line.</p>
</div>

<p align="center">
	<a href="https://github.com/eikendev/hackenv/actions"><img alt="Build status" src="https://img.shields.io/github/workflow/status/eikendev/hackenv/Main"/></a>&nbsp;
	<img alt="License" src="https://img.shields.io/github/license/eikendev/hackenv"/>&nbsp;
</p>

## 🚀&nbsp;Installation

```bash
go get -u github.com/eikendev/hackenv/cmd/...
```

Also check out the dependencies below.

## 🤘&nbsp;Features

- Download the latest official live image.
- Start an **SSH** daemon on the guest.
- Set up a **shared directory** between host and guest.
- Set the same **keyboard layout** in the guest as on the host.
- Configure the **firewall** to only allow SSH traffic from the host.

## 📄&nbsp;Usage

First, download an image using `hackenv get`.
This will download a live image from the official mirrors.
The download can take a while, so sit back and enjoy some tea.

By default, hackenv will operate with Kali Linux, and download its image.
If you want to work with Parrot Security instead, specify `hackenv --type=parrot`.

Next, run `hackenv up` to boot the virtual machine.
Once this command is finished, the VM is running and fully configured.

You can now start an SSH session with `hackenv ssh` or spin up a GUI with `hackenv gui`.

### File Sharing

hackenv will automatically set up a shared directory between the host and the virtual machine.
On the host side the directory is `~/.local/share/hackenv/shared`, while on the guest side it is located at `/shared`.

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
You can run `./bin/hackenv_fixlabels` if you are on Fedora or similar.
Be sure to re-adjust the permissions if you add files externally.

## ⚙&nbsp;Configuration

The tool currently does not support configuration via files.
However, some options can be set using environment variables.
For instance, to operate with Parrot Security by default, you can set `$HACKENV_TYPE=parrot`.
Check out the help (`--help`) of a command to see what other options support this.

## 🥙&nbsp;Dependencies

- [sshpass](https://sourceforge.net/projects/sshpass/)
- [libvirt](https://libvirt.org/)
- [virt-viewer](https://virt-manager.org/)

## 💡&nbsp;Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
