package commands

import (
	"crypto/sha256"
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
	slog.Info("Calculating checksum", "path", path)

	f, err := os.Open(path) //#nosec G304
	if err != nil {
		slog.Error("Failed to open file", "path", path, "err", err)
		return "", err
	}
	defer handling.Close(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		slog.Error("Failed to copy file content", "path", path, "err", err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// https://stackoverflow.com/a/11693049
func downloadImage(path, url string) error {
	slog.Info("Downloading image", "path", path, "url", url)

	out, err := os.Create(path) //#nosec G304
	if err != nil {
		slog.Error("Cannot create image file", "path", path, "err", err)
		return err
	}
	defer handling.Close(out)

	resp, err := http.Get(url) //#nosec G107
	if err != nil {
		slog.Error("Cannot download image file", "url", url, "err", err)
		return err
	}
	if resp == nil {
		return fmt.Errorf("received nil response")
	}
	defer handling.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		slog.Error("Cannot download image file: bad status", "status", resp.Status)
		return err
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		slog.Error("Cannot write image file", "path", path, "err", err)
		return err
	}

	slog.Info("Download successful")
	return nil
}

func validateChecksum(localPath, checksum string) error {
	newChecksum, err := calculateFileChecksum(localPath)
	if err != nil {
		return err
	}

	if newChecksum != checksum {
		err := os.Remove(localPath)
		if err != nil {
			slog.Error("Downloaded image has bad checksum and cannot be removed", "path", localPath, "expected", checksum, "actual", newChecksum, "err", err)
			os.Exit(1)
		}

		slog.Error("Downloaded image has bad checksum and was removed", "path", localPath, "expected", checksum, "actual", newChecksum)
		os.Exit(1)
	}

	slog.Info("Checksum validated successfully")
	return nil
}

// Run is the function for the get command.
func (c *GetCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)
	info := image.GetDownloadInfo(true)
	if info == nil {
		return fmt.Errorf("failed to get download information")
	}

	slog.Info("Found image to download", "filename", info.Filename, "checksum", info.Checksum)

	localPath := image.GetLocalPath(info.Version)

	// https://stackoverflow.com/a/12518877
	if _, err := os.Stat(localPath); err == nil {
		// The image already exists.

		if !c.Update && !c.Force {
			slog.Info("An image is already installed; use --update to refresh")
			return nil
		}

		localVersion := image.FileVersion(localPath)

		if !c.Force && image.VersionComparer.Eq(info.Version, localVersion) {
			slog.Info("Latest image is already installed; use --force to overwrite")
			return nil
		}
	} else if !os.IsNotExist(err) {
		slog.Error("Unable to get file information", "path", localPath, "err", err)
		os.Exit(1)
	}

	err := downloadImage(localPath, image.ArchiveURL+"/"+info.Filename)
	if err != nil {
		return err
	}

	err = validateChecksum(localPath, info.Checksum)
	if err != nil {
		return err
	}

	slog.Info("When using SELinux, label the image with the fix command before proceeding")

	return nil
}
