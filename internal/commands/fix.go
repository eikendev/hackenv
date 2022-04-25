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
	CreateBridge bool `short:"c" name:"create-bridge" help:"Create bridge"`
	RemoveBridge bool `short:"r" name:"remove-bridge" help:"Remove bridge"`
	Labels       bool `short:"l" name:"labels"        help:"Fix SElinux labels"`
}

// Run start the fix command
func (c *FixCommand) Run(s *options.Options) error {
	command := exec.Command("bash")

	if c.RemoveBridge {
		command.Stdin = strings.NewReader(scripts.RemoveBridgeScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
		return nil
	}

	if c.CreateBridge || (!c.CreateBridge && !c.Labels) {
		command.Stdin = strings.NewReader(scripts.CreateBridgeScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}

	if c.RemoveBridge {
		command.Stdin = strings.NewReader(scripts.RemoveBridgeScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}

	return nil
}
