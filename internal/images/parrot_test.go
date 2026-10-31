package images

import (
	"fmt"
	"regexp"
	"runtime"
	"testing"
)

const sampleParrotIndex = `<html><body><pre>
<a href="../">../</a>
<a href="4.11.3/">4.11.3/</a>
<a href="5.2-1/">5.2-1/</a>
<a href="6.3.2/">6.3.2/</a>
<a href="7.10/">7.10/</a>
<a href="7.4/">7.4/</a>
<a href="caine/">caine/</a>
<a href="latest/">latest/</a>
</pre></body></html>`

func TestFindLatestRelease(t *testing.T) {
	got, err := findLatestRelease(sampleParrotIndex, getGenericVersionComparer())
	if err != nil {
		t.Fatalf("findLatestRelease() error = %v", err)
	}
	if got != "7.10" {
		t.Fatalf("findLatestRelease() = %q, want %q", got, "7.10")
	}
}

func TestFindLatestReleaseWithoutReleases(t *testing.T) {
	_, err := findLatestRelease(`<a href="latest/">latest/</a>`, getGenericVersionComparer())
	if err == nil {
		t.Fatal("findLatestRelease() succeeded on an index without releases")
	}
}

func sampleParrotHashes(withTuned bool) []byte {
	tuned := ""
	if withTuned {
		tuned = fmt.Sprintf("bbbb  Parrot-security-tuned-7.4_%s.iso\n", runtime.GOARCH)
	}

	return fmt.Appendf(nil, `-----BEGIN PGP SIGNED MESSAGE-----
Hash: SHA512

md5
1111  Parrot-security-7.4_%[1]s.iso

sha256
aaaa  Parrot-home-7.4_%[1]s.iso
cccc  Parrot-security-7.4_%[1]s.iso
%[2]sdddd  Parrot-security-7.4_%[1]s.iso.torrent
eeee  Parrot-security-7.4_rpi.img.xz

sha512
2222  Parrot-security-7.4_%[1]s.iso
`, runtime.GOARCH, tuned)
}

func TestParrotDownloadInfoPrefersTunedImage(t *testing.T) {
	versionRegex := regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

	info, err := parrotDownloadInfo(sampleParrotHashes(true), versionRegex)
	if err != nil {
		t.Fatalf("parrotDownloadInfo() error = %v", err)
	}

	wantFile := fmt.Sprintf("Parrot-security-tuned-7.4_%s.iso", runtime.GOARCH)
	if info.Filename != wantFile || info.Checksum != "bbbb" || info.Version != "7.4" {
		t.Fatalf("parrotDownloadInfo() = %+v, want %s with checksum bbbb", info, wantFile)
	}
}

func TestParrotDownloadInfoFallsBackToPlainImage(t *testing.T) {
	versionRegex := regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

	info, err := parrotDownloadInfo(sampleParrotHashes(false), versionRegex)
	if err != nil {
		t.Fatalf("parrotDownloadInfo() error = %v", err)
	}

	if info.Filename != fmt.Sprintf("Parrot-security-7.4_%s.iso", runtime.GOARCH) || info.Checksum != "cccc" {
		t.Fatalf("parrotDownloadInfo() = %+v, want plain image with checksum cccc", info)
	}
}

func TestParrotDownloadInfoWithoutSecurityImage(t *testing.T) {
	versionRegex := regexp.MustCompile(`\d+\.\d+`)

	sums := fmt.Appendf(nil, "sha256\naaaa  Parrot-home-7.4_%s.iso\n", runtime.GOARCH)
	if _, err := parrotDownloadInfo(sums, versionRegex); err == nil {
		t.Fatal("parrotDownloadInfo() succeeded without a security image")
	}
}
