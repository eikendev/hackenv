// Package images provides utilities to access image information and manage images.
package images

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/adrg/xdg"
	rawLibvirt "libvirt.org/go/libvirt"

	"github.com/eikendev/hackenv/internal/constants"
)

type infoRetriever func(string, *regexp.Regexp) (*DownloadInfo, error)

type bootInitializer func(*rawLibvirt.Domain)

type sshStarter func(*rawLibvirt.Domain)

type versionComparer interface {
	Lt(string, string) bool
	Gt(string, string) bool
	Eq(string, string) bool
}

// Image contains information about a stored image.
type Image struct {
	Name              string
	DisplayName       string
	ArchiveURL        string
	checksumPath      string
	LocalImageName    string
	VersionRegex      *regexp.Regexp
	SSHUser           string
	SSHPassword       string
	MacAddress        string
	infoRetriever     infoRetriever
	bootInitializer   bootInitializer
	sshStarter        sshStarter
	ConfigurationCmds []string
	VersionComparer   versionComparer
}

// DownloadInfo contains information about an image that can be downloaded.
type DownloadInfo struct {
	Checksum string
	Version  string
	Filename string
}

var images = map[string]Image{
	"kali": {
		Name:              "kali",
		DisplayName:       "Kali Linux",
		ArchiveURL:        "https://kali.download/base-images/current/",
		checksumPath:      "/SHA256SUMS",
		LocalImageName:    "kali-%s.iso",
		VersionRegex:      regexp.MustCompile(`\d\d\d\d\.\d+`),
		SSHUser:           "kali",
		SSHPassword:       "kali",
		MacAddress:        "52:54:00:08:f9:e8",
		infoRetriever:     kaliInfoRetriever,
		bootInitializer:   kaliBootInitializer,
		sshStarter:        kaliSSHStarter,
		ConfigurationCmds: kaliConfigurationCmds,
		VersionComparer:   getGenericVersionComparer(),
	},
	"parrot": {
		Name:              "parrot",
		DisplayName:       "Parrot Security",
		ArchiveURL:        "https://deb.parrot.sh/parrot/iso/current",
		checksumPath:      "/signed-hashes.txt",
		LocalImageName:    "parrot-%s.iso",
		VersionRegex:      regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`),
		SSHUser:           "user",
		SSHPassword:       "parrot",
		MacAddress:        "52:54:00:08:f9:e9",
		infoRetriever:     parrotInfoRetriever,
		bootInitializer:   parrotBootInitializer,
		sshStarter:        parrotSSHStarter,
		ConfigurationCmds: []string{},
		VersionComparer:   getGenericVersionComparer(),
	},
}

// GetDownloadInfo retreives the necessary information to download an image.
func (i *Image) GetDownloadInfo(strict bool) *DownloadInfo {
	info, err := i.infoRetriever(i.ArchiveURL+i.checksumPath, i.VersionRegex)
	if err != nil && strict {
		slog.Error("Cannot retrieve latest image details", "image", i.DisplayName, "err", err)
		os.Exit(1)
	}

	return info
}

// Boot executes the necessary steps to boot a downloaded image.
func (i *Image) Boot(dom *rawLibvirt.Domain, version string) {
	slog.Info("Booting image", "image", i.DisplayName, "version", version)
	i.bootInitializer(dom)
}

// StartSSH executes the necessary steps to start SSH on a booted image.
func (i *Image) StartSSH(dom *rawLibvirt.Domain) {
	slog.Info("Bootstrapping SSH", "image", i.DisplayName)
	i.sshStarter(dom)
}

// GetLocalPath builds the full path of a downloaded image based on a given version.
func (i *Image) GetLocalPath(version string) string {
	filename := fmt.Sprintf(i.LocalImageName, version)

	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, filename))
	if err != nil {
		slog.Error("Cannot access data directory", "err", err, "file", filename)
		os.Exit(1)
	}

	return path
}

// GetLatestPath returns the full path of the image with the greatest version.
func (i *Image) GetLatestPath() string {
	imageGlob := fmt.Sprintf(i.LocalImageName, "*")

	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, imageGlob))
	if err != nil {
		slog.Error("Cannot access data directory", "err", err, "pattern", imageGlob)
		os.Exit(1)
	}

	matches, err := filepath.Glob(path)
	if err != nil || len(matches) == 0 {
		slog.Error("Cannot find image", "image", i.DisplayName, "err", err)
		os.Exit(1)
		return "" // Won't actually return due to log.Fatal
	}

	latestPath := matches[0]
	latestVersion := i.FileVersion(latestPath)

	for _, path := range matches {
		slog.Info("Found image path", "path", path)
		if newVersion := i.FileVersion(path); i.VersionComparer.Gt(newVersion, latestVersion) {
			latestPath = path
			latestVersion = newVersion
		}
	}

	return latestPath
}

// GetImageDetails returns detailed information about a given image.
func GetImageDetails(name string) Image {
	image, ok := images[name]
	if !ok {
		slog.Error("Image not supported", "image", name)
		os.Exit(1)
	}
	return image
}

// GetAllImages returns a map of all available images.
func GetAllImages() map[string]Image {
	return images
}

// FileVersion returns the version of the image given its full path.
func (i *Image) FileVersion(path string) string {
	return i.VersionRegex.FindString(path)
}

func sendKeys(dom *rawLibvirt.Domain, keys []uint) {
	err := dom.SendKey(uint(rawLibvirt.KEYCODE_SET_LINUX), 10, keys, 0)
	if err != nil {
		slog.Error("Cannot send keys", "err", err)
		os.Exit(1)
	}

	time.Sleep(20 * time.Millisecond)
}
