package commands

import (
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"

	"github.com/eikendev/hackenv/internal/doctor"
	"github.com/eikendev/hackenv/internal/options"
)

// DoctorCommand represents the options specific to the doctor command.
type DoctorCommand struct{}

var statusColors = map[doctor.Status]*color.Color{
	doctor.Pass: color.New(color.FgGreen),
	doctor.Fail: color.New(color.FgRed),
	doctor.Skip: color.New(color.FgYellow),
}

// Run is the function for the doctor command.
func (*DoctorCommand) Run(_ *options.Options) error {
	results := doctor.RunAll(doctor.Checks())

	writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, result := range results {
		_, _ = fmt.Fprintf(writer, "%s\t%s\t%s\n", statusColors[result.Status].Sprint(result.Status), result.Name, result.Detail)
		if result.Hint != "" {
			_, _ = fmt.Fprintf(writer, "\t\t-> %s\n", result.Hint)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("cannot write results: %w", err)
	}

	if doctor.HasFailures(results) {
		return errors.New("some checks failed")
	}

	return nil
}
