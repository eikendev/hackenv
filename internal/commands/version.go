package commands

import (
	"fmt"

	"github.com/eikendev/hackenv/internal/buildconfig"
	"github.com/eikendev/hackenv/internal/options"
)

type VersionCommand struct{}

func (c *VersionCommand) Run(s *options.Options) error {
	fmt.Printf("hackenv %s\n", buildconfig.Version)

	return nil
}
