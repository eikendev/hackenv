package commands

import (
	_ "embed"
	"fmt"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/options"
)

// FixCommand struct
type FixCommand struct {
	CreateBridge bool `short:"c" name:"create-bridge" help:"Create bridge"`
	RemoveBridge bool `short:"r" name:"remove-bridge" help:"Remove bridge"`
	Labels       bool `short:"l" name:"labels" help:"Fix selinux labels"`
}

//go:embed hackenv_createbridge
var createBridgeScript string

//go:embed hackenv_removebridge
var removeBridgeScript string

//go:embed hackenv_fixlabels
var fixLabelsScript string

// Run start the fix command
func (c *FixCommand) Run(s *options.Options) error {
	command := exec.Command("bash")

	if c.CreateBridge || (!c.CreateBridge && !c.Labels) {
		command.Stdin = strings.NewReader(createBridgeScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}

	if c.Labels || (!c.CreateBridge && !c.Labels) {
		command.Stdin = strings.NewReader(fixLabelsScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}

	if c.RemoveBridge {
		command.Stdin = strings.NewReader(removeBridgeScript)
		b, err := command.Output()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(b))
	}

	return nil
}
