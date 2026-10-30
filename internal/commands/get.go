package commands

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	progressbar "github.com/schollz/progressbar/v3"

	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/network"
	"github.com/eikendev/hackenv/internal/options"
)

const (
	partialSuffix = ".part"
	selinuxMount  = "/sys/fs/selinux"
)

// GetCommand represents the options specific to the get command.
type GetCommand struct {
	Force         bool   `short:"f" name:"force" help:"Force to download the new image"`
	Update        bool   `short:"u" name:"update" help:"Allow update to the latest image"`
	SkipSignature bool   `name:"skip-signature" help:"Do not verify the signature of the checksum file"`
	SigningKey    string `name:"signing-key" env:"HACKENV_SIGNING_KEY" type:"existingfile" help:"Armored public key to verify the checksum file with instead of the embedded one"`
}

func (c *GetCommand) fetchOptions() (images.FetchOptions, error) {
	opts := images.FetchOptions{SkipSignature: c.SkipSignature}

	if c.SigningKey != "" {
		key, err := os.ReadFile(c.SigningKey) //#nosec G304
		if err != nil {
			slog.Error("failed to read signing key", "path", c.SigningKey, "err", err)
			return opts, fmt.Errorf("cannot read signing key %s: %w", c.SigningKey, err)
		}
		opts.SigningKey = key
	}

	return opts, nil
}

// interruptContext returns a context that is cancelled on the first interrupt.
// A second interrupt terminates the process as usual.
func interruptContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	go func() {
		<-ctx.Done()
		stop()
	}()

	return ctx, stop
}

// https://golang.org/pkg/crypto/sha256/#example_New_file
func calculateFileChecksum(path string) (string, error) {
	slog.Info("calculating checksum", "path", path)

	f, err := os.Open(path) //#nosec G304
	if err != nil {
		slog.Error("failed to open file", "path", path, "err", err)
		return "", fmt.Errorf("cannot open %s: %w", path, err)
	}
	defer handling.Close(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		slog.Error("failed to read file", "path", path, "err", err)
		return "", fmt.Errorf("cannot read %s for checksum: %w", path, err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func downloadToFile(ctx context.Context, path, url string) error {
	out, err := os.Create(path) //#nosec G304
	if err != nil {
		slog.Error("failed to create image file", "path", path, "err", err)
		return fmt.Errorf("cannot create %s: %w", path, err)
	}
	defer handling.Close(out)

	resp, err := network.GetResponse(ctx, url)
	if err != nil {
		slog.Error("failed to download image", "url", url, "err", err)
		return fmt.Errorf("cannot download %s: %w", url, err)
	}
	defer handling.Close(resp.Body)

	bar := progressbar.DefaultBytes(resp.ContentLength, "downloading")

	if _, err := io.Copy(io.MultiWriter(out, bar), resp.Body); err != nil {
		slog.Error("failed to write image file", "path", path, "err", err)
		return fmt.Errorf("cannot write %s: %w", path, err)
	}

	return nil
}

// downloadImage downloads url into a partial file next to path and renames it into place.
func downloadImage(ctx context.Context, path, url string) error {
	slog.Info("downloading image", "path", path, "url", url)

	partial := path + partialSuffix

	if err := downloadToFile(ctx, partial, url); err != nil {
		removePartial(partial)
		return err
	}

	if err := os.Rename(partial, path); err != nil {
		removePartial(partial)
		slog.Error("failed to move image into place", "from", partial, "to", path, "err", err)
		return fmt.Errorf("cannot move %s to %s: %w", partial, path, err)
	}

	slog.Info("downloaded image", "path", path)
	return nil
}

func removePartial(path string) {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		slog.Warn("failed to remove partial download", "path", path, "err", err)
	}
}

func validateChecksum(localPath, checksum string) error {
	newChecksum, err := calculateFileChecksum(localPath)
	if err != nil {
		slog.Error("failed to calculate checksum", "path", localPath, "err", err)
		return fmt.Errorf("failed to calculate checksum for %s: %w", localPath, err)
	}

	if newChecksum != checksum {
		removeErr := os.Remove(localPath)
		if removeErr != nil {
			slog.Error("failed to remove invalid image", "path", localPath, "expected", checksum, "actual", newChecksum, "err", removeErr)
			return fmt.Errorf("failed to remove image with bad checksum (expected %s, actual %s): %w", checksum, newChecksum, removeErr)
		}

		slog.Error("found invalid checksum", "path", localPath, "expected", checksum, "actual", newChecksum)
		return fmt.Errorf("detected bad checksum for downloaded image: expected %s, actual %s", checksum, newChecksum)
	}

	slog.Info("validated checksum")
	return nil
}

func hintSELinuxLabels() {
	if _, err := os.Stat(selinuxMount); err != nil {
		return
	}

	slog.Warn("missing selinux labels on new image", "hint", "run hackenv fix apply-labels")
}

// Run is the function for the get command.
func (c *GetCommand) Run(s *options.Options) error {
	ctx, stop := interruptContext()
	defer stop()

	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("failed to resolve image", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}
	opts, err := c.fetchOptions()
	if err != nil {
		return err
	}

	info, err := image.GetDownloadInfo(ctx, true, opts)
	if err != nil {
		slog.Error("failed to fetch download info", "image", image.DisplayName, "err", err)
		return fmt.Errorf("cannot fetch download info for %s: %w", image.DisplayName, err)
	}
	if info == nil {
		slog.Error("missing download info", "image", image.DisplayName)
		return fmt.Errorf("failed to get download information")
	}

	slog.Info("found image to download", "filename", info.Filename, "checksum", info.Checksum)

	localPath, err := image.GetLocalPath(info.Version)
	if err != nil {
		slog.Error("failed to resolve image path", "image", image.DisplayName, "version", info.Version, "err", err)
		return fmt.Errorf("cannot resolve local path for %s %s: %w", image.DisplayName, info.Version, err)
	}

	skipDownload, err := c.shouldSkipDownload(localPath, image, info)
	if err != nil {
		slog.Error("failed to stat image", "path", localPath, "err", err)
		return fmt.Errorf("cannot decide whether to download %s: %w", localPath, err)
	}
	if skipDownload {
		return nil
	}

	if err := downloadImage(ctx, localPath, info.URL); err != nil {
		return fmt.Errorf("cannot download image %s: %w", info.URL, err)
	}

	if err := validateChecksum(localPath, info.Checksum); err != nil {
		return fmt.Errorf("failed to validate checksum for %s: %w", localPath, err)
	}

	hintSELinuxLabels()

	return nil
}

func (c *GetCommand) shouldSkipDownload(localPath string, image images.Image, info *images.DownloadInfo) (bool, error) {
	// https://stackoverflow.com/a/12518877
	_, err := os.Stat(localPath)
	if err == nil {
		if !c.Update && !c.Force {
			slog.Info("image already exists", "hint", "use --update to refresh")
			return true, nil
		}

		localVersion := image.FileVersion(localPath)

		if !c.Force && image.VersionComparer.Eq(info.Version, localVersion) {
			slog.Info("already up to date", "hint", "use --force to overwrite")
			return true, nil
		}
		return false, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	slog.Error("failed to stat image", "path", localPath, "err", err)
	return false, fmt.Errorf("cannot stat %s: %w", localPath, err)
}
