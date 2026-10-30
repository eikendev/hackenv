package commands

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	progressbar "github.com/schollz/progressbar/v3"

	"github.com/eikendev/hackenv/internal/handling"
	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/options"
)

// GetCommand represents the options specific to the get command.
type GetCommand struct {
	Force  bool `short:"f" name:"force" help:"Force to download the new image"`
	Update bool `short:"u" name:"update" help:"Allow update to the latest image"`
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

// https://stackoverflow.com/a/11693049
func downloadImage(path, url string) error {
	slog.Info("downloading image", "path", path, "url", url)

	out, err := os.Create(path) //#nosec G304
	if err != nil {
		slog.Error("failed to create image file", "path", path, "err", err)
		return fmt.Errorf("cannot create %s: %w", path, err)
	}
	defer handling.Close(out)

	resp, err := http.Get(url) //#nosec G107
	if err != nil {
		slog.Error("failed to download image", "url", url, "err", err)
		return fmt.Errorf("cannot download %s: %w", url, err)
	}
	if resp == nil {
		return fmt.Errorf("received nil response")
	}
	defer handling.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		slog.Error("failed to download image", "status", resp.Status)
		return fmt.Errorf("cannot download %s: status %s", url, resp.Status)
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		slog.Error("failed to write image file", "path", path, "err", err)
		return fmt.Errorf("cannot write %s: %w", path, err)
	}

	slog.Info("downloaded image")
	return nil
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

// Run is the function for the get command.
func (c *GetCommand) Run(s *options.Options) error {
	image, err := images.GetImageDetails(s.Type)
	if err != nil {
		slog.Error("failed to resolve image", "type", s.Type, "err", err)
		return fmt.Errorf("cannot resolve image details for %q: %w", s.Type, err)
	}
	info, err := image.GetDownloadInfo(true)
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

	err = downloadImage(localPath, image.ArchiveURL+"/"+info.Filename)
	if err != nil {
		slog.Error("failed to download image", "path", localPath, "url", image.ArchiveURL+"/"+info.Filename, "err", err)
		return fmt.Errorf("cannot download image %s/%s: %w", image.ArchiveURL, info.Filename, err)
	}

	err = validateChecksum(localPath, info.Checksum)
	if err != nil {
		slog.Error("failed to validate checksum", "path", localPath, "err", err)
		return fmt.Errorf("failed to validate checksum for %s: %w", localPath, err)
	}

	slog.Info("skipping selinux labels", "hint", "run hackenv fix apply-labels when using SELinux")

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
