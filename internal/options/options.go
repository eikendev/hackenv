// Package options defines the global options of this tool.
package options

import "github.com/alecthomas/kong"

// Options represents the global options of this tool.
type Options struct {
	Config    kong.ConfigFlag `name:"config" env:"HACKENV_CONFIG" help:"Load options from this YAML file instead of the default one"`
	Type      string          `name:"type" env:"HACKENV_TYPE" default:"kali" enum:"kali,parrot" help:"The VM to control with this command"`
	Provision bool            `name:"provision" env:"HACKENV_PROVISION" help:"Provision the VM"`
	Verbose   bool            `name:"verbose" short:"v" help:"Verbose mode"`
}

// Runnable defines an interface for subcommands that take the global settings and a password.
type Runnable interface {
	Run(*Options) error
}
