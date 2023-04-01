// Package commands contains functions that are exposed as dedicated commands of the tool.
package commands

import (
	"fmt"

	"github.com/eikendev/hackenv/internal/buildconfig"
	"github.com/eikendev/hackenv/internal/options"
)

// VersionCommand represents the options specific to the version command.
type VersionCommand struct{}

// Run is the function for the version command.
func (*VersionCommand) Run(_ *options.Options) error {
	fmt.Printf("hackenv %s\n", buildconfig.Version)

	return nil
}
