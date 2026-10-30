package host

import "testing"

func TestParseLayout(t *testing.T) {
	tests := []struct {
		name    string
		output  string
		key     string
		want    string
		wantErr bool
	}{
		{
			name:   "setxkbmap",
			output: "rules:      evdev\nmodel:      pc105\nlayout:     at,us\nvariant:    nodeadkeys,\n",
			key:    "layout",
			want:   "at,us",
		},
		{
			name:   "localectl",
			output: "System Locale: LANG=en_US.UTF-8\n    VC Keymap: at-nodeadkeys\n   X11 Layout: at,us\n    X11 Model: pc105\n",
			key:    "x11 layout",
			want:   "at,us",
		},
		{
			name:    "missing key",
			output:  "rules:      evdev\nmodel:      pc105\n",
			key:     "layout",
			wantErr: true,
		},
		{
			name:    "empty value",
			output:  "layout:\n",
			key:     "layout",
			wantErr: true,
		},
		{
			name:    "no output",
			output:  "",
			key:     "layout",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseLayout(tt.output, tt.key)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseLayout() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("parseLayout() = %q, want %q", got, tt.want)
			}
		})
	}
}
