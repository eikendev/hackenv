package images

import (
	"fmt"
	"regexp"
	"strings"
)

// parseChecksumLine extracts checksum information from a line
func parseChecksumLine(line string, versionRegex *regexp.Regexp) (*DownloadInfo, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid checksum line format: %s", line)
	}

	checksum := parts[0]
	filename := parts[len(parts)-1]

	return &DownloadInfo{
		Checksum: checksum,
		Version:  versionRegex.FindString(filename),
		Filename: filename,
	}, nil
}
