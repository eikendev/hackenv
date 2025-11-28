// Package main provides the main function as a starting point of this tool.
package main

import (
	"github.com/alecthomas/kong"

	"github.com/eikendev/hackenv/internal/commands"
	"github.com/eikendev/hackenv/internal/logging"
	"github.com/eikendev/hackenv/internal/options"
)

var cmd struct {
	options.Options
	Down    commands.DownCommand    `cmd:"down" aliases:"d" help:"Shut down the VM"`
	Get     commands.GetCommand     `cmd:"get" help:"Download the VM image"`
	GUI     commands.GuiCommand     `cmd:"gui" aliases:"g" help:"Open a GUI for the VM"`
	SSH     commands.SSHCommand     `cmd:"ssh" aliases:"s" help:"Open an SSH session for the VM"`
	Status  commands.StatusCommand  `cmd:"status" help:"Print the status of the VM"`
	Up      commands.UpCommand      `cmd:"up" aliases:"u" help:"Initialize and start the VM or provision if already running"`
	Fix     commands.FixCommand     `cmd:"fix" aliases:"f" help:"Fix helpers: manage bridge and apply SELinux labels"`
	Version commands.VersionCommand `cmd:"version" help:"Print the version of hackenv"`
}

func main() {
	kctx := kong.Parse(&cmd,
		kong.Description("hackenv provisions and manages preconfigured VMs for security research."),
		kong.UsageOnError(),
		kong.Bind(&cmd.Options),
	)

	logging.Setup(cmd.Verbose)

	err := kctx.Run()
	kctx.FatalIfErrorf(err)
}
