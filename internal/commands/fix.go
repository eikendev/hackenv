package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/scripts"
)

// FixCommand struct
type FixCommand struct {
	CreateBridge createBridge `cmd:"create-bridge" aliases:"c" help:"Create bridge"`
	RemoveBridge removeBridge `cmd:"remove-bridge" aliases:"r" help:"Remove bridge"`
	FixLabels    fixLabels    `cmd:"fix-labels" aliases:"l" help:"Fix SElinux labels"`
	All          all          `cmd:"all" aliases:"a" help:"Create bridge and fix selinux labels"`
}

type createBridge struct{}

func (c *createBridge) Run(s *options.Options) error {
	return command(scripts.CreateBridgeScript, s.Verbose)
}

type removeBridge struct{}

func (c *removeBridge) Run(s *options.Options) error {
	return command(scripts.RemoveBridgeScript, s.Verbose)
}

type fixLabels struct{}

func (c *fixLabels) Run(s *options.Options) error {
	return command(scripts.FixLabelsScript, s.Verbose)
}

type all struct{}

func (c *all) Run(s *options.Options) error {
	err := command(scripts.CreateBridgeScript, s.Verbose)
	if err != nil {
		return err
	}
	return command(scripts.FixLabelsScript, s.Verbose)
}

func command(script string, verbose bool) error {
	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(script)
	b, err := cmd.Output()
	if verbose {
		fmt.Println(string(b))
	}
	if err != nil {
		return err
	}

	return nil

}
