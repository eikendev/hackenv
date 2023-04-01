// Package constants contains global constants.
package constants

const (
	// XdgAppname is the name of this tool for XDG purposes.
	XdgAppname = "hackenv"

	// ConnectURI is the URI where the QEMU instance is made available.
	ConnectURI = "qemu:///session"

	// SSHKeypairName is the name of the SSH keypair used to connect to the guest.
	SSHKeypairName = "sshkey"

	// PostbootFile is the filename of the script used to populate the guest.
	PostbootFile = "postboot.sh"
)
