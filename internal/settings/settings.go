package settings

type Settings struct {
	//lint:ignore SA5008 go-flags makes use of duplicate struct tags
	Type string `long:"type" env:"HACKENV_TYPE" default:"kali" choice:"kali" choice:"parrot" description:"The VM to control with this command"`
}

// Runnable defines an interface for subcommands that take the global settings and a password.
type Runnable interface {
	Run(*Settings)
}

// Runner is the subcommand to run after all arguments were parsed.
var Runner Runnable
