package commands

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/eikendev/hackenv/internal/images"
	"github.com/eikendev/hackenv/internal/settings"
)

type GetCommand struct {
	Force  bool `short:"f" long:"force" description:"Force to download the new image"`
	Update bool `short:"u" long:"update" description:"Allow update to the latest image"`
}

func (c *GetCommand) Execute(args []string) error {
	settings.Runner = c
	return nil
}

// https://golang.org/pkg/crypto/sha256/#example_New_file
func calculateFileChecksum(path string) string {
	log.Printf("Calculating checksum of %s\n", path)

	f, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %s\n", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatalf("Failed to copy file content: %s\n", err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// https://stackoverflow.com/a/11693049
func downloadImage(path string, url string) {
	log.Printf("Downloading image to %s\n", path)

	out, err := os.Create(path)
	if err != nil {
		log.Fatalf("Cannot write image file: %s\n", err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Cannot download image file: %s\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Cannot download image file: bad status %s\n", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("Cannot write image file: %s\n", err)
	}

	log.Println("Download successful")
}

func (c *GetCommand) Run(s *settings.Settings) {
	image := images.GetImageDetails(s.Type)

	checksum, filename := image.GetDownloadInfo()

	log.Printf("Found file %s with checksum %s\n", filename, checksum)

	localPath := image.GetLocalPath()

	// https://stackoverflow.com/a/12518877
	if _, err := os.Stat(localPath); !c.Force && !os.IsNotExist(err) {
		// The image already exists.

		if !c.Update {
			log.Println("An image is already installed; update with --update")
			return
		}

		if checksum == calculateFileChecksum(localPath) {
			log.Println("Latest image of VM already install; force with --force")
			return
		}
	}

	downloadImage(localPath, image.ArchiveURL+"/"+filename)
}
