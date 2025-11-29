package commands

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/options"
	"github.com/eikendev/hackenv/internal/scripts"
)

// FixCommand represents the options specific to the fix command.
type FixCommand struct {
	CreateBridge createBridge `cmd:"create-bridge" aliases:"c" help:"Create bridge"`
	RemoveBridge removeBridge `cmd:"remove-bridge" aliases:"r" help:"Remove bridge"`
	ApplyLabels  applyLabels  `cmd:"apply-labels" aliases:"l" help:"Apply SELinux labels"`
	All          all          `cmd:"all" aliases:"a" help:"Create bridge and apply SELinux labels"`
}

type createBridge struct{}

// Run is the function for the run command.
func (c *createBridge) Run(s *options.Options) error {
	return execCommand([]string{scripts.CreateBridgeScript}, s.Verbose)
}

type removeBridge struct{}

func (c *removeBridge) Run(s *options.Options) error {
	return execCommand([]string{scripts.RemoveBridgeScript}, s.Verbose)
}

type applyLabels struct{}

func (c *applyLabels) Run(s *options.Options) error {
	return execCommand([]string{scripts.ApplyLabelsScript}, s.Verbose)
}

type all struct{}

func (c *all) Run(s *options.Options) error {
	return execCommand([]string{scripts.CreateBridgeScript, scripts.ApplyLabelsScript}, s.Verbose)
}

func execCommand(scripts []string, verbose bool) error {
	for i, script := range scripts {
		cmd := exec.Command("bash")
		slog.Info("Running script", "index", i+1, "total", len(scripts))
		cmd.Stdin = strings.NewReader(script)
		b, err := cmd.CombinedOutput()
		if verbose {
			fmt.Println(string(b))
		}
		if err != nil {
			slog.Error("Script execution failed", "index", i+1, "total", len(scripts), "err", err)
			return fmt.Errorf("failed to run script %d/%d: %w", i+1, len(scripts), err)
		}
	}

	return nil
}
