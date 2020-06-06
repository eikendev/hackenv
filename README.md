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
You can either do this manually, or by instrumenting `make download` to download the latest release from the official mirrors.

If downloading the image manually, make sure to use the live version.
You have to put the file into this directory, or create a symbolic link, and rename it to `kali.iso`.

A `make install` creates (and boots) the virtual machine (domain); `make gui` will open a graphical user interface.
You can enable SSH inside the box with `sudo service ssh restart`, and connect to it from your host using `make ssh`.

### File Sharing

A mounting hint for the `shared/` directory in this repository is created as `/shared`.
It can be mounted from inside the guest using the following snippet with privileges.
See the [KVM documentation](https://www.linux-kvm.org/page/9p_virtio) for further information.
```bash
mkdir /shared
mount -t 9p -o trans=virtio,version=9p2000.L /shared /shared
```


If you have trouble with SELinux, you have to adjust the context of the `shared/` directory.
The make target `make permissions` will do this for you if you are running Fedora or similar.
