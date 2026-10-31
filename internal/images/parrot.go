package images

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"time"

	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/network"
)

var parrotReleaseDirRegex = regexp.MustCompile(`href="(\d+(?:\.\d+)+)/"`)

var (
	parrotSecurityImageRegex = regexp.MustCompile(`^Parrot-security(?:-[a-z]+)*-\d+(?:\.\d+)*_` + runtime.GOARCH + `\.iso$`)
	parrotTunedImageRegex    = regexp.MustCompile(`^Parrot-security-tuned-\d`)
	parrotPlainImageRegex    = regexp.MustCompile(`^Parrot-security-\d`)
)

func findParrotChecksums(scanner *bufio.Scanner) (map[string]string, error) {
	var inSha256Section bool

	checksums := map[string]string{}

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		switch {
		case len(fields) == 1:
			inSha256Section = fields[0] == "sha256"
		case len(fields) == 2 && inSha256Section && parrotSecurityImageRegex.MatchString(fields[1]):
			checksums[fields[1]] = fields[0]
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("cannot read checksum file: %w", err)
	}

	if len(checksums) == 0 {
		return nil, errors.New("cannot find checksum in file")
	}

	return checksums, nil
}

func pickParrotImage(checksums map[string]string) string {
	for _, preferred := range []*regexp.Regexp{parrotTunedImageRegex, parrotPlainImageRegex} {
		for _, filename := range slices.Sorted(maps.Keys(checksums)) {
			if preferred.MatchString(filename) {
				return filename
			}
		}
	}

	return ""
}

func parrotDownloadInfo(sums []byte, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	checksums, err := findParrotChecksums(bufio.NewScanner(bytes.NewReader(sums)))
	if err != nil {
		return nil, err
	}

	filename := pickParrotImage(checksums)
	if filename == "" {
		return nil, errors.New("cannot find security image in checksum file")
	}

	return &DownloadInfo{
		Checksum: checksums[filename],
		Version:  versionRegex.FindString(filename),
		Filename: filename,
	}, nil
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

	info, err := parrotDownloadInfo(sums, image.VersionRegex)
	if err != nil {
		slog.Error("missing checksum line", "err", err)
		return nil, fmt.Errorf("cannot find Parrot checksum line: %w", err)
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
