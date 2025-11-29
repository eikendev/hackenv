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

var kaliConfigurationCmds = []string{
	"touch ~/.hushlogin",
}

func findKaliChecksumLine(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "live-"+runtime.GOARCH+".iso") {
			return strings.TrimSpace(line), nil
		}
	}

	return "", errors.New("cannot find checksum in file")
}

func kaliInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := network.GetResponse(url)
	if err != nil {
		slog.Error("Failed to fetch Kali checksum file", "url", url, "err", err)
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("Failed to close response body", "err", err)
		}
	}()

	line, err := findKaliChecksumLine(bufio.NewScanner(resp.Body))
	if err != nil {
		slog.Error("Failed to find Kali checksum line", "err", err)
		return nil, fmt.Errorf("cannot find Kali checksum line: %w", err)
	}

	return parseChecksumLine(line, versionRegex)
}

func kaliBootInitializer(dom *rawLibvirt.Domain) error {
	if err := genericBootInitializer(dom); err != nil {
		slog.Error("Failed to boot Kali image", "err", err)
		return fmt.Errorf("failed to boot Kali image: %w", err)
	}
	return nil
}

func kaliSSHStarter(dom *rawLibvirt.Domain) error {
	time.Sleep(5 * time.Second)
	if err := switchToTTY(dom); err != nil {
		slog.Error("Failed to switch Kali console to TTY", "err", err)
		return fmt.Errorf("failed to switch Kali console to TTY: %w", err)
	}

	time.Sleep(1 * time.Second)
	if err := systemdRestartSSH(dom); err != nil {
		slog.Error("Failed to restart SSH on Kali image", "err", err)
		return fmt.Errorf("failed to restart SSH on Kali image: %w", err)
	}

	if err := switchFromTTY(dom); err != nil {
		slog.Error("Failed to switch Kali console back from TTY", "err", err)
		return fmt.Errorf("failed to switch Kali console from TTY: %w", err)
	}

	return nil
}
