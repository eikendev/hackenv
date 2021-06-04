package images

import (
	"log"
	"path/filepath"
	"time"

	"github.com/adrg/xdg"
	"github.com/eikendev/hackenv/internal/constants"
	rawLibvirt "libvirt.org/libvirt-go"
)

type infoRetriever func(string) (string, string)

type bootInitializer func(*rawLibvirt.Domain)

type sshStarter func(*rawLibvirt.Domain)

type Image struct {
	Name              string
	DisplayName       string
	ArchiveURL        string
	checksumPath      string
	LocalImageName    string
	SSHUser           string
	SSHPassword       string
	MacAddress        string
	infoRetriever     infoRetriever
	bootInitializer   bootInitializer
	sshStarter        sshStarter
	ConfigurationCmds []string
}

var images = map[string]Image{
	"kali": {
		Name:              "kali",
		DisplayName:       "Kali Linux",
		ArchiveURL:        "https://cdimage.kali.org/current",
		checksumPath:      "/SHA256SUMS",
		LocalImageName:    "kali.iso",
		SSHUser:           "kali",
		SSHPassword:       "kali",
		MacAddress:        "52:54:00:08:f9:e8",
		infoRetriever:     kaliInfoRetriever,
		bootInitializer:   kaliBootInitializer,
		sshStarter:        kaliSSHStarter,
		ConfigurationCmds: kaliConfigurationCmds,
	},
	"parrot": {
		Name:              "parrot",
		DisplayName:       "Parrot Security",
		ArchiveURL:        "https://download.parrot.sh/parrot/iso/current",
		checksumPath:      "/signed-hashes.txt",
		LocalImageName:    "parrot.iso",
		SSHUser:           "user",
		SSHPassword:       "toor",
		MacAddress:        "52:54:00:08:f9:e9",
		infoRetriever:     parrotInfoRetriever,
		bootInitializer:   parrotBootInitializer,
		sshStarter:        parrotSSHStarter,
		ConfigurationCmds: []string{},
	},
}

func (i *Image) GetDownloadInfo() (string, string) {
	return i.infoRetriever(i.ArchiveURL + "/" + i.checksumPath)
}

func (i *Image) Boot(dom *rawLibvirt.Domain) {
	log.Printf("Booting %s\n", i.DisplayName)
	i.bootInitializer(dom)
}

func (i *Image) StartSSH(dom *rawLibvirt.Domain) {
	log.Printf("Starting SSH daemon\n")
	i.sshStarter(dom)
}

func (i *Image) GetLocalPath() string {
	path, err := xdg.DataFile(filepath.Join(constants.XdgAppname, i.LocalImageName))
	if err != nil {
		log.Fatalf("Cannot access data directory: %s\n", err)
	}
	return path
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

func sendKeys(dom *rawLibvirt.Domain, keys []uint) {
	dom.SendKey(uint(rawLibvirt.KEYCODE_SET_LINUX), 10, keys, 0)
	time.Sleep(20 * time.Millisecond)
}
