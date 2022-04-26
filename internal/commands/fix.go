package commands

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/scripts"
	log "github.com/sirupsen/logrus"
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
	return execCommand([]string{scripts.CreateBridgeScript}, s.Verbose)
}

type removeBridge struct{}

func (c *removeBridge) Run(s *options.Options) error {
	return execCommand([]string{scripts.RemoveBridgeScript}, s.Verbose)
}

type fixLabels struct{}

func (c *fixLabels) Run(s *options.Options) error {
	return execCommand([]string{scripts.FixLabelsScript}, s.Verbose)
}

type all struct{}

func (c *all) Run(s *options.Options) error {
	if err := execCommand([]string{scripts.CreateBridgeScript, scripts.FixLabelsScript}, s.Verbose); err != nil {
		return err
	}

	return nil
}

func execCommand(scripts []string, verbose bool) error {

	for i, script := range scripts {
		cmd := exec.Command("bash")
		log.Infof("Running script %v/%v...\n\n", i, len(scripts)-1)
		cmd.Stdin = strings.NewReader(script)
		b, err := cmd.CombinedOutput()
		if verbose {
			fmt.Println(string(b))
		}
		if err != nil {
			return err
		}
	}

	return nil
}
