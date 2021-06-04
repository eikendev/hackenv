package images

import (
	"bufio"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"

	rawLibvirt "libvirt.org/libvirt-go"
)

var kaliConfigurationCmds = []string{
	"touch ~/.hushlogin",
}

func kaliInfoRetriever(url string) (string, string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Cannot retrieve latest image: %s\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Cannot retrieve latest image: bad status %s\n", resp.Status)
	}

	var line string
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "live") && strings.Contains(line, runtime.GOARCH) {
			break
		}
	}

	if line == "" {
		log.Fatalf("Cannot retrieve latest image: bad checksum file\n")
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	checksum := parts[0]
	filename := parts[len(parts)-1]

	return checksum, filename
}

func kaliBootInitializer(dom *rawLibvirt.Domain) {
	genericBootInitializer(dom)
}

func kaliSSHStarter(dom *rawLibvirt.Domain) {
	time.Sleep(2 * time.Second)

	switchToTTY(dom)
	systemdRestartSSH(dom)
	switchFromTTY(dom)
}
