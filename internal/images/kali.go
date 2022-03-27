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

var kaliConfigurationCmds = []string{
	"touch ~/.hushlogin",
}

func kaliInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad HTTP status code (%s)", resp.Status)
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

func kaliBootInitializer(dom *rawLibvirt.Domain) {
	genericBootInitializer(dom)
}

func kaliSSHStarter(dom *rawLibvirt.Domain) {
	time.Sleep(5 * time.Second)
	switchToTTY(dom)
	time.Sleep(1 * time.Second)
	systemdRestartSSH(dom)
	switchFromTTY(dom)
}
