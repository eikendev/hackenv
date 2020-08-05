[![Build status](https://img.shields.io/travis/eikendev/kali-libvirt/master)](https://travis-ci.com/github/eikendev/kali-libvirt/builds/)
![License](https://img.shields.io/github/license/eikendev/kali-libvirt)

## About

Did you ever find yourself in a situation where you needed access to a fresh [Kali Linux](https://www.kali.org/) installation and didn't want to go through the trouble of setting it up manually?
Naturally, you can either use
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/) to automate setting up a virtual machine, or
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/) to utilize Linux namespaces (containers).

Personally, I sometimes _want_ to use a virtual machine to enforce stricter separation.
However, [libvirt](https://libvirt.org/) turned out to be a lot more convenient to use on my system.
Since I could not find a Kali Linux image providing libvirt support for Vagrant, I settled to use this helper script.

## Usage

![Showcase](https://i.imgur.com/YxIJvCz.gif)

To make the following steps more convenient, I advise you to add `./bin/kali` to your path.
For instance, you can run `make install` to link it from your `~/bin` directory.

First, download the Kali Linux image that you want to use.
You can either do this manually, or by instrumenting `kali download` to download the latest release from the official mirrors.

If downloading the image manually, make sure to use the live version.
You have to put the file into this directory or create a symbolic link, and rename it to `kali.iso`.

To get started, use
- `kali install` to create and boot the virtual machine (domain), and
- `kali gui` to open a graphical user interface.

### SSH Access

Once the machine has booted, enable SSH inside the box with `sudo service ssh restart`.
At this point, you can connect to it from your host using `kali ssh`.

### File Sharing

You can now run `kali share` to setup a shared directory between the host and the virtual machine.
On the host side, the directory is `./shared` inside this repository.
On the client side, it is located at `/shared`.

If SELinux denies access to the shared directory, you have to adjust the context of the directory.
Running `kali permissions` will do this for you if you are on Fedora or similar.
Be sure to re-adjust the permissions if you add files externally.

## Dependencies

This should be a complete list.
Most of these are standard tools.

- [curl](https://curl.haxx.se/)
- [jq](https://stedolan.github.io/jq/)
- [sshpass](https://sourceforge.net/projects/sshpass/)
- [virsh](https://libvirt.org/)
- [virt-install](https://virt-manager.org/)
