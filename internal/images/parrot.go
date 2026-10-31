package images

import (
	"bufio"
	"bytes"
	"context"
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

func latestParrotRelease(ctx context.Context, image *Image) (string, error) {
	index, err := network.GetBytes(ctx, image.ArchiveURL+"/", maxIndexSize)
	if err != nil {
		slog.Error("failed to fetch release index", "url", image.ArchiveURL, "err", err)
		return "", fmt.Errorf("failed to fetch release index: %w", err)
	}

	return findLatestRelease(string(index), image.VersionComparer)
}

func fetchParrotChecksums(ctx context.Context, image *Image, releaseURL string, opts FetchOptions) ([]byte, error) {
	message, err := network.GetBytes(ctx, releaseURL+image.checksumPath, maxChecksumFileSize)
	if err != nil {
		slog.Error("failed to fetch checksum file", "url", releaseURL+image.checksumPath, "err", err)
		return nil, fmt.Errorf("failed to fetch checksum file: %w", err)
	}

	if opts.SkipSignature {
		return message, nil
	}

	return verifyClearsigned(image.signingKeyFor(opts), message)
}

func parrotInfoRetriever(ctx context.Context, image *Image, opts FetchOptions) (*DownloadInfo, error) {
	release, err := latestParrotRelease(ctx, image)
	if err != nil {
		return nil, err
	}

	releaseURL := image.ArchiveURL + "/" + release

	sums, err := fetchParrotChecksums(ctx, image, releaseURL, opts)
	if err != nil {
		return nil, err
	}

	line, err := findParrotChecksumLine(bufio.NewScanner(bytes.NewReader(sums)))
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

	if err := switchToDesktop(dom, KEY_F8); err != nil {
		slog.Error("failed to switch from tty", "err", err)
		return fmt.Errorf("failed to switch Parrot console back from TTY: %w", err)
	}

	return nil
}
