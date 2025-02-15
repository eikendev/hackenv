package commands

import (
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"os"

	progressbar "github.com/schollz/progressbar/v3"
	log "github.com/sirupsen/logrus"

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
	log.Printf("Calculating checksum of %s\n", path)

	f, err := os.Open(path) //#nosec G304
	if err != nil {
		log.Errorf("Failed to open file: %s\n", err)
		return "", err
	}
	defer handling.Close(f)

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Errorf("Failed to copy file content: %s\n", err)
		return "", err
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// https://stackoverflow.com/a/11693049
func downloadImage(path, url string) error {
	log.Printf("Downloading image to %s\n", path)

	out, err := os.Create(path) //#nosec G304
	if err != nil {
		log.Errorf("Cannot write image file: %s\n", err)
		return err
	}
	defer handling.Close(out)

	resp, err := http.Get(url) //#nosec G107
	if err != nil {
		log.Errorf("Cannot download image file: %s\n", err)
		return err
	}
	if resp == nil {
		return fmt.Errorf("received nil response")
	}
	defer handling.Close(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Cannot download image file: bad status %s\n", resp.Status)
		return err
	}

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading",
	)

	_, err = io.Copy(io.MultiWriter(out, bar), resp.Body)
	if err != nil {
		log.Errorf("Cannot write image file: %s\n", err)
		return err
	}

	log.Println("Download successful")
	return nil
}

func validateChecksum(localPath, checksum string) error {
	newChecksum, err := calculateFileChecksum(localPath)
	if err != nil {
		return err
	}

	if newChecksum != checksum {
		checksumMsg := fmt.Sprintf("Downloaded image has bad checksum: %s instead of %s", newChecksum, checksum)

		err := os.Remove(localPath)
		if err != nil {
			log.Fatalf("%s. Unable to remove file.\n", checksumMsg)
		}

		log.Fatalf("%s. File removed.\n", checksumMsg)
	}

	log.Println("Checksum validated successfully")
	return nil
}

// Run is the function for the get command.
func (c *GetCommand) Run(s *options.Options) error {
	image := images.GetImageDetails(s.Type)
	info := image.GetDownloadInfo(true)
	if info == nil {
		return fmt.Errorf("failed to get download information")
	}

	log.Printf("Found file %s with checksum %s\n", info.Filename, info.Checksum)

	localPath := image.GetLocalPath(info.Version)

	// https://stackoverflow.com/a/12518877
	if _, err := os.Stat(localPath); err == nil {
		// The image already exists.

		if !c.Update && !c.Force {
			log.Println("An image is already installed; update with --update")
			return nil
		}

		localVersion := image.FileVersion(localPath)

		if !c.Force && image.VersionComparer.Eq(info.Version, localVersion) {
			log.Println("Latest image is already installed; force with --force")
			return nil
		}
	} else if !os.IsNotExist(err) {
		log.Fatalf("Unable to get file information for path %s\n", localPath)
	}

	err := downloadImage(localPath, image.ArchiveURL+"/"+info.Filename)
	if err != nil {
		return err
	}

	err = validateChecksum(localPath, info.Checksum)
	if err != nil {
		return err
	}

	log.Info("When using SELinux, don't forget to label the image with the fix command before proceeding")

	return nil
}
