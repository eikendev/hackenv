package images

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"runtime"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	rawLibvirt "libvirt.org/libvirt-go"

	"github.com/eikendev/hackenv/internal/network"
)

var kaliConfigurationCmds = []string{
	"touch ~/.hushlogin",
}

func findKaliChecksumLine(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "live-"+runtime.GOARCH+".iso") {
			return strings.TrimSpace(line), nil
		}
	}

	return "", errors.New("checksum not found in file")
}

func kaliInfoRetriever(url string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	resp, err := network.GetResponse(url)
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Warnf("failed to close response body: %v", err)
		}
	}()

	line, err := findKaliChecksumLine(bufio.NewScanner(resp.Body))
	if err != nil {
		return nil, err
	}

	return parseChecksumLine(line, versionRegex)
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
