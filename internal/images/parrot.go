package images

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"regexp"
	"runtime"
	"strings"
	"time"

	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/network"
)

func findParrotChecksumLine(scanner *bufio.Scanner) (string, error) {
	const (
		sha256Section = "sha256"
		sha384Section = "sha384"
	)

	var inSha256Section bool

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.Contains(line, sha256Section):
			inSha256Section = true
			continue
		case strings.Contains(line, sha384Section):
			inSha256Section = false
			continue
		case inSha256Section &&
			strings.Contains(line, "Parrot-security") &&
			strings.Contains(line, "_"+runtime.GOARCH+".iso"):
			return strings.TrimSpace(line), nil
		}
	}

	return "", errors.New("cannot find checksum in file")
}

func parrotInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := network.GetResponse(url)
	if err != nil {
		slog.Error("failed to fetch checksum file", "url", url, "err", err)
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "err", err)
		}
	}()

	line, err := findParrotChecksumLine(bufio.NewScanner(resp.Body))
	if err != nil {
		slog.Error("missing checksum line", "err", err)
		return nil, fmt.Errorf("cannot find Parrot checksum line: %w", err)
	}

	return parseChecksumLine(line, versionRegex)
}

func parrotBootInitializer(dom *rawLibvirt.Domain) error {
	if err := genericBootInitializer(dom); err != nil {
		slog.Error("failed to boot image", "err", err)
		return fmt.Errorf("failed to boot Parrot image: %w", err)
	}
	return nil
}

func parrotSSHStarter(dom *rawLibvirt.Domain) error {
	time.Sleep(5 * time.Second)
	if err := switchToTTY(dom); err != nil {
		slog.Error("failed to switch to tty", "err", err)
		return fmt.Errorf("failed to switch Parrot console to TTY: %w", err)
	}

	time.Sleep(1 * time.Second)
	if err := enablePasswordSSH(dom); err != nil {
		slog.Error("failed to enable password ssh", "err", err)
		return fmt.Errorf("failed to enable password SSH on Parrot image: %w", err)
	}

	if err := systemdRestartSSH(dom); err != nil {
		slog.Error("failed to restart ssh", "err", err)
		return fmt.Errorf("failed to restart SSH on Parrot image: %w", err)
	}

	if err := switchFromTTY(dom); err != nil {
		slog.Error("failed to switch from tty", "err", err)
		return fmt.Errorf("failed to switch Parrot console from TTY: %w", err)
	}

	return nil
}
