// Package host provides various utilities related to the host.
package host

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

const meminfoPath = "/proc/meminfo"

// TotalMemoryKiB returns the physical memory of the host in KiB.
func TotalMemoryKiB() (int64, error) {
	file, err := os.Open(meminfoPath)
	if err != nil {
		return 0, fmt.Errorf("cannot open %s: %w", meminfoPath, err)
	}
	defer func() { _ = file.Close() }()

	return parseTotalMemory(file)
}

func parseTotalMemory(r io.Reader) (int64, error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		name, value, found := strings.Cut(scanner.Text(), ":")
		if !found || name != "MemTotal" {
			continue
		}

		fields := strings.Fields(value)
		if len(fields) < 1 {
			break
		}

		total, err := strconv.ParseInt(fields[0], 10, 64)
		if err != nil || total <= 0 {
			return 0, fmt.Errorf("found invalid MemTotal %q", value)
		}

		return total, nil
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("cannot read meminfo: %w", err)
	}

	return 0, errors.New("missing MemTotal in meminfo")
}
