package commands

import (
	"fmt"
	"log/slog"
	"os"
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
func (c *applyLabels) Run(_ *options.Options) error {
	return execScript(scripts.ApplyLabelsScript)
}

func execScript(script string) error {
	cmd := exec.Command("bash")
	cmd.Stdin = strings.NewReader(script)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	slog.Info("running script")

	if err := cmd.Run(); err != nil {
		slog.Error("failed to run script", "err", err)
		return fmt.Errorf("failed to run script: %w", err)
	}

	return nil
}
