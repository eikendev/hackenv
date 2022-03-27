package images

import (
	"bufio"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"runtime"
	"strings"
	"time"

	rawLibvirt "libvirt.org/libvirt-go"
)

func parrotInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTTP status code (%s)", resp.Status)
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
		return nil, errors.New("bad checksum file")
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	filename := parts[len(parts)-1]

	return &DownloadInfo{
		parts[0],
		versionRegex.FindString(filename),
		filename,
	}, nil
}

func parrotBootInitializer(dom *rawLibvirt.Domain) {
	genericBootInitializer(dom)
}

func parrotSSHStarter(dom *rawLibvirt.Domain) {
	time.Sleep(5 * time.Second)
	switchToTTY(dom)
	time.Sleep(1 * time.Second)
	enablePasswordSSH(dom)
	systemdRestartSSH(dom)
	switchFromTTY(dom)
}
