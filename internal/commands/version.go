// Package commands contains functions that are exposed as dedicated commands of the tool.
package commands

import (
	"fmt"
	"runtime/debug"

	"github.com/eikendev/hackenv/internal/options"
)

// VersionCommand represents the options specific to the version command.
type VersionCommand struct{}

// Run is the function for the version command.
func (*VersionCommand) Run(_ *options.Options) error {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("build info not available")
	}

	fmt.Printf("hackenv %s\n", buildInfo.Main.Version)
	return nil
}
