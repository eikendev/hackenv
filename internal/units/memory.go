// Package units parses human readable sizes used on the command line.
package units

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Memory is an amount of RAM in KiB. Flags parse it from values such as 4G or 4096M.
type Memory int64

const (
	kib = 1
	mib = 1024 * kib
	gib = 1024 * mib
)

var unitMultipliers = map[string]int64{
	"":    mib,
	"k":   kib,
	"kib": kib,
	"m":   mib,
	"mib": mib,
	"g":   gib,
	"gib": gib,
}

// ParseMemory converts a size such as "4G", "4096M", "2097152K" or a bare number of MiB to KiB.
func ParseMemory(value string) (Memory, error) {
	text := strings.ToLower(strings.TrimSpace(value))
	if text == "" {
		return 0, errors.New("memory size is empty")
	}

	digits := strings.TrimRight(text, "abcdefghijklmnopqrstuvwxyz")
	unit := strings.TrimSpace(text[len(digits):])

	multiplier, ok := unitMultipliers[unit]
	if !ok {
		return 0, fmt.Errorf("unknown memory unit %q in %q", unit, value)
	}

	amount, err := strconv.ParseInt(strings.TrimSpace(digits), 10, 64)
	if err != nil || amount <= 0 {
		return 0, fmt.Errorf("invalid memory size %q", value)
	}

	return Memory(amount * multiplier), nil
}

// UnmarshalText implements encoding.TextUnmarshaler so kong can parse flags.
func (m *Memory) UnmarshalText(text []byte) error {
	parsed, err := ParseMemory(string(text))
	if err != nil {
		return err
	}

	*m = parsed
	return nil
}

// KiB returns the size in KiB.
func (m Memory) KiB() int64 {
	return int64(m)
}

// String formats the size with the largest exact unit.
func (m Memory) String() string {
	switch {
	case m%gib == 0:
		return fmt.Sprintf("%dG", m/gib)
	case m%mib == 0:
		return fmt.Sprintf("%dM", m/mib)
	default:
		return fmt.Sprintf("%dK", m)
	}
}
