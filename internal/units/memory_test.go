package units

import "testing"

func TestParseMemory(t *testing.T) {
	tests := []struct {
		input   string
		want    Memory
		wantErr bool
	}{
		{input: "4G", want: 4 * gib},
		{input: "4GiB", want: 4 * gib},
		{input: " 4 g ", want: 4 * gib},
		{input: "4096M", want: 4096 * mib},
		{input: "4096", want: 4096 * mib},
		{input: "2097152K", want: 2097152 * kib},
		{input: "0", wantErr: true},
		{input: "-1G", wantErr: true},
		{input: "4T", wantErr: true},
		{input: "lots", wantErr: true},
		{input: "", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got, err := ParseMemory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseMemory(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ParseMemory(%q) = %d, want %d", tt.input, got, tt.want)
			}
		})
	}
}

func TestMemoryString(t *testing.T) {
	tests := map[Memory]string{
		4 * gib:    "4G",
		1536 * mib: "1536M",
		1000 * kib: "1000K",
	}

	for m, want := range tests {
		if got := m.String(); got != want {
			t.Errorf("Memory(%d).String() = %q, want %q", m, got, want)
		}
	}
}

func TestMemoryUnmarshalText(t *testing.T) {
	var m Memory
	if err := m.UnmarshalText([]byte("2G")); err != nil || m != 2*gib {
		t.Fatalf("UnmarshalText() = %d, %v", m, err)
	}
	if err := m.UnmarshalText([]byte("nope")); err == nil {
		t.Fatal("UnmarshalText() accepted garbage")
	}
}
