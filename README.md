[![Build status](https://img.shields.io/github/workflow/status/eikendev/hackenv/Main)](https://github.com/eikendev/hackenv/actions)
![License](https://img.shields.io/github/license/eikendev/hackenv)

# hackenv

## About

This tool is for you if you frequently use [Kali Linux](https://www.kali.org/) or [Parrot Security](https://www.parrotsec.org/).
It makes it easy to retrieve an up-to-date image and get it running via [libvirt](https://libvirt.org/).

## Installation

The following command will download, build, and install the tool.

```bash
go get -u github.com/eikendev/hackenv/cmd/...
```

## Usage

First, you need to download an image using `hackenv get`.
This will download a live image from the official mirrors.
The download can take a while, so sit back and enjoy some tea.

Next, run `hackenv up` to boot the virtual machine.
Once this command is finished, the VM is running and fully configured.

You can now decide to start an SSH session with `hackenv ssh` or spin up a GUI with `hackenv gui`.

### File Sharing

hackenv will automatically try to setup a shared directory between the host and the virtual machine.
On the host side, the directory is `~/.local/share/hackenv/shared`.
On the guest side, it is located at `/shared`.

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
You can run `./bin/hackenv_fixlabels` if you are on Fedora or similar.
Be sure to re-adjust the permissions if you add files externally.

## Configuration

The tool currently does not support configuration via files.
However, some options can be set using environment variables.
Check out the help (`--help`) to see what options support this.

## Dependencies

- [sshpass](https://sourceforge.net/projects/sshpass/)
- [virsh](https://libvirt.org/)
- [virt-viewer](https://virt-manager.org/)

## Alternatives

If you do not like this tool, the following options are worth checking out:
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/)
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/)
