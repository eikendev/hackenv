package host

import (
	"strings"
	"testing"
)

func TestParseTotalMemory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{name: "fedora", input: "MemTotal:       32764584 kB\nMemFree:        16000000 kB\n", want: 32764584},
		{name: "missing", input: "MemFree:        16000000 kB\n", wantErr: true},
		{name: "garbage", input: "MemTotal:       lots kB\n", wantErr: true},
		{name: "empty", input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseTotalMemory(strings.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseTotalMemory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("parseTotalMemory() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTotalMemoryKiB(t *testing.T) {
	total, err := TotalMemoryKiB()
	if err != nil {
		t.Fatalf("TotalMemoryKiB() error = %v", err)
	}
	if total <= 0 {
		t.Fatalf("TotalMemoryKiB() = %d", total)
	}
}
