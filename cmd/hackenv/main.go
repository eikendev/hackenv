package main

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/eikendev/hackenv/internal/commands"
	"github.com/eikendev/hackenv/internal/settings"
	"github.com/jessevdk/go-flags"
)

type command struct {
	settings.Settings
	Down   commands.DownCommand   `command:"down" alias:"d" description:"Shut down the VM"`
	Get    commands.GetCommand    `command:"get" description:"Download the VM image"`
	GUI    commands.GuiCommand    `command:"gui" alias:"g" description:"Open a GUI for the VM"`
	SSH    commands.SSHCommand    `command:"ssh" alias:"s" description:"Open an SSH session for the VM"`
	Status commands.StatusCommand `command:"status" description:"Print the status of the VM"`
	Up     commands.UpCommand     `command:"up" alias:"u" description:"Initialize and start the VM or provision if already running"`
}

var (
	cmds   command
	parser = flags.NewParser(&cmds, flags.Default)
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableTimestamp: true,
	})

	log.SetOutput(os.Stdout)

	log.SetLevel(log.InfoLevel)
}

func main() {
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	settings.Runner.Run(&cmds.Settings)
}
