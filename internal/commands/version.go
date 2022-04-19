package commands

import (
	"fmt"

	"github.com/eikendev/hackenv/internal/buildconfig"
	"github.com/eikendev/hackenv/internal/settings"
)

type VersionCommand struct{}

func (c *VersionCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

func (c *VersionCommand) Run(s *settings.Settings) {
	fmt.Printf("hackenv %s\n", buildconfig.Version)
}
