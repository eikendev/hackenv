// Package logging provides utilities for configuring and managing application logging.
package logging

import (
	"log/slog"
	"os"
)

// Setup configures the global slog logger based on the verbose flag.
func Setup(verbose bool) {
	logLevel := slog.LevelInfo
	if verbose {
		logLevel = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})
	slog.SetDefault(slog.New(handler))

	if verbose {
		slog.Info("Verbose logging enabled", "level", logLevel)
	}
}
