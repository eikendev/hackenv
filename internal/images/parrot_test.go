package images

import "testing"

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
