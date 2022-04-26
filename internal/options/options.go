package options

type Options struct {
	Type      string `name:"type" env:"HACKENV_TYPE" default:"kali" enum:"kali,parrot" help:"The VM to control with this command"`
	Keymap    string `name:"keymap" env:"HACKENV_KEYMAP" default:"" help:"The keyboard keymap to force"`
	Provision bool   `name:"provision" env:"HACKENV_PROVISION" help:"Provision the VM"`
	Verbose   bool   `name:"verbose" short:"v" help:"Verbose mode"`
}

// Runnable defines an interface for subcommands that take the global settings and a password.
type Runnable interface {
	Run(*Options) error
}
