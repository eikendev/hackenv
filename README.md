## About

Did you ever find yourself in a situation where you needed access to a fresh [Kali Linux](https://www.kali.org/) installation and didn't want to go through the trouble of setting it up manually?
Naturally, you can either use
- [Vagrant](https://www.vagrantup.com/) in combination with [VirtualBox](https://www.virtualbox.org/) to automate setting up a virtual machine, or
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/) to utilize Linux namespaces (containers).

Personally, I sometimes _want_ to use a virtual machine to enforce stricter separation.
However, [libvirt](https://libvirt.org/) turned out to be a lot more convenient to use on my system.
Since I could not find a Kali Linux image providing libvirt support for Vagrant, I settled to use this helper script.

## Usage

First, download the Kali Linux image that you want to use.
This script assumes that it is the live version, no hard drive will be configured.
Either put the file into this directory, or create a symbolic link, and rename it as `kali.iso`.

A `make install` creates the virtual machine (domain), with `make gui` will open a graphical user interface.
You can enable SSH inside the box with `service ssh restart`, and connect to it from your host using `make ssh`.
