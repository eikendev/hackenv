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
	ApplyLabels applyLabels `cmd:"apply-labels" aliases:"l" help:"Apply SELinux labels to images and the shared directory"`
}

type applyLabels struct{}

// Run is the function for the apply-labels command.
func (c *applyLabels) Run(s *options.Options) error {
	return execScript(scripts.ApplyLabelsScript, s.Verbose)
}

func execScript(script string, verbose bool) error {
	cmd := exec.Command("bash")
	slog.Info("running script")
	cmd.Stdin = strings.NewReader(script)
	b, err := cmd.CombinedOutput()
	if verbose {
		fmt.Println(string(b))
	}
	if err != nil {
		slog.Error("failed to run script", "err", err)
		return fmt.Errorf("failed to run script: %w", err)
	}

	return nil
}
