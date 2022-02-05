package images

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/adrg/xdg"
	"github.com/eikendev/hackenv/internal/constants"
	rawLibvirt "libvirt.org/libvirt-go"
)

type infoRetriever func(string, *regexp.Regexp) (*DownloadInfo, error)

type bootInitializer func(*rawLibvirt.Domain)

type sshStarter func(*rawLibvirt.Domain)

type versionComparer interface {
	Lt(string, string) bool
	Gt(string, string) bool
	Eq(string, string) bool
}

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

type DownloadInfo struct {
	Checksum string
	Version  string
	Filename string
}

var images = map[string]Image{
	"kali": {
		Name:              "kali",
		DisplayName:       "Kali Linux",
		ArchiveURL:        "https://cdimage.kali.org/current",
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
		ArchiveURL:        "https://download.parrot.sh/parrot/iso/5.0",
		checksumPath:      "/signed-hashes.txt",
		LocalImageName:    "parrot-%s.iso",
		VersionRegex:      regexp.MustCompile(`\d+\.\d+`),
		SSHUser:           "user",
		SSHPassword:       "toor",
		MacAddress:        "52:54:00:08:f9:e9",
		infoRetriever:     parrotInfoRetriever,
		bootInitializer:   parrotBootInitializer,
		sshStarter:        parrotSSHStarter,
		ConfigurationCmds: []string{},
		VersionComparer:   getGenericVersionComparer(),
	},
}

func (i *Image) GetDownloadInfo(strict bool) *DownloadInfo {
	info, err := i.infoRetriever(i.ArchiveURL+i.checksumPath, i.VersionRegex)
	if err != nil && strict {
		log.Fatalf("Cannot retrieve latest image details: %s\n", err)
	}

	return info
}

func (i *Image) Boot(dom *rawLibvirt.Domain, version string) {
	log.Printf("Booting %s %s\n", i.DisplayName, version)
	i.bootInitializer(dom)
}

func (i *Image) StartSSH(dom *rawLibvirt.Domain) {
	log.Printf("Bootstraping...\n")
	i.sshStarter(dom)
}

func (i *Image) GetLocalPath(version string) string {
	filename := fmt.Sprintf(i.LocalImageName, version)

	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, filename))
	if err != nil {
		log.Fatalf("Cannot access data directory: %s\n", err)
	}

	return path
}

func (i *Image) GetLatestPath() string {
	imageGlob := fmt.Sprintf(i.LocalImageName, "*")

	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, imageGlob))
	if err != nil {
		log.Fatalf("Cannot access data directory: %s\n", err)
	}

	matches, err := filepath.Glob(path)
	if err != nil {
		log.Fatalf("Malformed glob pattern: %s\n", err)
	}

	if matches == nil {
		log.Fatalf("Image for %s not found; download with get command\n", i.DisplayName)
	}

	latestPath := matches[0]
	latestVersion := i.FileVersion(latestPath)

	for _, path := range matches {
		log.Printf("Found image path %s\n", path)
		if newVersion := i.FileVersion(path); i.VersionComparer.Gt(newVersion, latestVersion) {
			latestPath = path
			latestVersion = newVersion
		}
	}

	return latestPath
}

func GetImageDetails(name string) Image {
	image, ok := images[name]
	if !ok {
		log.Fatalf("Image not supported: %s", name)
	}
	return image
}

func GetAllImages() map[string]Image {
	return images
}

func (i *Image) FileVersion(path string) string {
	return i.VersionRegex.FindString(path)
}

func sendKeys(dom *rawLibvirt.Domain, keys []uint) {
	dom.SendKey(uint(rawLibvirt.KEYCODE_SET_LINUX), 10, keys, 0)
	time.Sleep(20 * time.Millisecond)
}
