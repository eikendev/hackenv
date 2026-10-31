package images

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
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

func fetchKaliChecksums(ctx context.Context, image *Image, opts FetchOptions) ([]byte, error) {
	sums, err := network.GetBytes(ctx, image.ArchiveURL+image.checksumPath, maxChecksumFileSize)
	if err != nil {
		slog.Error("failed to fetch checksum file", "url", image.ArchiveURL+image.checksumPath, "err", err)
		return nil, fmt.Errorf("failed to fetch checksum file: %w", err)
	}

	if opts.SkipSignature {
		return sums, nil
	}

	signature, err := network.GetBytes(ctx, image.ArchiveURL+image.signaturePath, maxSignatureSize)
	if err != nil {
		slog.Error("failed to fetch signature file", "url", image.ArchiveURL+image.signaturePath, "err", err)
		return nil, fmt.Errorf("failed to fetch signature file: %w", err)
	}

	return verifyDetached(image.signingKeyFor(opts), sums, signature)
}

func kaliInfoRetriever(ctx context.Context, image *Image, opts FetchOptions) (*DownloadInfo, error) {
	sums, err := fetchKaliChecksums(ctx, image, opts)
	if err != nil {
		return nil, err
	}

	line, err := findKaliChecksumLine(bufio.NewScanner(bytes.NewReader(sums)))
	if err != nil {
		slog.Error("missing checksum line", "err", err)
		return nil, fmt.Errorf("cannot find Kali checksum line: %w", err)
	}

	info, err := parseChecksumLine(line, image.VersionRegex)
	if err != nil {
		return nil, err
	}

	info.URL = image.ArchiveURL + "/" + info.Filename
	info.TorrentURL = info.URL + ".torrent"
	return info, nil
}

func kaliBootInitializer(dom *rawLibvirt.Domain) error {
	if err := genericBootInitializer(dom); err != nil {
		slog.Error("failed to boot image", "err", err)
		return fmt.Errorf("failed to boot Kali image: %w", err)
	}
	return nil
}

func kaliSSHStarter(dom *rawLibvirt.Domain) error {
	time.Sleep(5 * time.Second)
	if err := switchToTTY(dom); err != nil {
		slog.Error("failed to switch to tty", "err", err)
		return fmt.Errorf("failed to switch Kali console to TTY: %w", err)
	}

	time.Sleep(1 * time.Second)
	if err := systemdRestartSSH(dom); err != nil {
		slog.Error("failed to restart ssh", "err", err)
		return fmt.Errorf("failed to restart SSH on Kali image: %w", err)
	}

	if err := switchToDesktop(dom, KEY_F7); err != nil {
		slog.Error("failed to switch from tty", "err", err)
		return fmt.Errorf("failed to switch Kali console back from TTY: %w", err)
	}

	return nil
}
