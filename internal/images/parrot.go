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

	return "", errors.New("checksum not found in file")
}

func parrotInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := network.GetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("Failed to close response body", "err", err)
		}
	}()

	line, err := findParrotChecksumLine(bufio.NewScanner(resp.Body))
	if err != nil {
		return nil, err
	}

	return parseChecksumLine(line, versionRegex)
}

func parrotBootInitializer(dom *rawLibvirt.Domain) {
	genericBootInitializer(dom)
}

func parrotSSHStarter(dom *rawLibvirt.Domain) {
	time.Sleep(5 * time.Second)
	switchToTTY(dom)
	time.Sleep(1 * time.Second)
	enablePasswordSSH(dom)
	systemdRestartSSH(dom)
	switchFromTTY(dom)
}
