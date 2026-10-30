package images

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"regexp"
	"runtime"
	"strings"
	"time"

	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/network"
)

const maxIndexSize = 1 << 20

var parrotReleaseDirRegex = regexp.MustCompile(`href="(\d+(?:\.\d+)+)/"`)

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
			strings.Contains(line, "Parrot-security-") &&
			strings.Contains(line, "_"+runtime.GOARCH+".iso"):
			return strings.TrimSpace(line), nil
		}
	}

	return "", errors.New("cannot find checksum in file")
}

func findLatestRelease(index string, comparer versionComparer) (string, error) {
	var latest string

	for _, match := range parrotReleaseDirRegex.FindAllStringSubmatch(index, -1) {
		version := match[1]
		if latest == "" || comparer.Gt(version, latest) {
			latest = version
		}
	}

	if latest == "" {
		return "", errors.New("cannot find release directories in index")
	}

	return latest, nil
}

func latestParrotRelease(image *Image) (string, error) {
	resp, err := network.GetResponse(image.ArchiveURL + "/")
	if err != nil {
		slog.Error("failed to fetch release index", "url", image.ArchiveURL, "err", err)
		return "", fmt.Errorf("failed to fetch release index: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			slog.Warn("failed to close response body", "err", err)
		}
	}()

	index, err := io.ReadAll(io.LimitReader(resp.Body, maxIndexSize))
	if err != nil {
		slog.Error("failed to read release index", "url", image.ArchiveURL, "err", err)
		return "", fmt.Errorf("failed to read release index: %w", err)
	}

	return findLatestRelease(string(index), image.VersionComparer)
}

func parrotInfoRetriever(image *Image) (*DownloadInfo, error) {
	release, err := latestParrotRelease(image)
	if err != nil {
		return nil, err
	}

	releaseURL := image.ArchiveURL + "/" + release

	resp, err := network.GetResponse(releaseURL + image.checksumPath)
	if err != nil {
		slog.Error("failed to fetch checksum file", "url", releaseURL+image.checksumPath, "err", err)
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

	info, err := parseChecksumLine(line, image.VersionRegex)
	if err != nil {
		return nil, err
	}

	info.URL = releaseURL + "/" + info.Filename
	return info, nil
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
		return fmt.Errorf("failed to switch Parrot console back from TTY: %w", err)
	}

	return nil
}
