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

func parrotInfoRetriever(url string) (string, string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Cannot retrieve latest image: %s\n", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Cannot retrieve latest image: bad status %s\n", resp.Status)
	}

	var line string
	skipped := false
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		line = scanner.Text()

		if skipped {
			if strings.Contains(line, "t-security") && strings.Contains(line, runtime.GOARCH) {
				break
			}
		} else {
			if strings.Contains(line, "sha256") {
				skipped = true
			}
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

func parrotBootInitializer(dom *rawLibvirt.Domain) {
	genericBootInitializer(dom)
}

func parrotSSHStarter(dom *rawLibvirt.Domain) {
	time.Sleep(2 * time.Second)

	switchToTTY(dom)
	enablePasswordSSH(dom)
	systemdRestartSSH(dom)
	switchFromTTY(dom)
}
