package guest

import (
	"net"
	"strconv"
	"testing"
	"time"
)

func TestParseSSHPort(t *testing.T) {
	tests := []struct {
		name    string
		desc    string
		want    int
		wantErr bool
	}{
		{
			name: "forwarded ssh port",
			desc: `<domain><devices><interface type='user'>
				<portForward proto='tcp' address='127.0.0.1'><range start='2222' to='22'/></portForward>
			</interface></devices></domain>`,
			want: 2222,
		},
		{
			name: "udp forward is ignored",
			desc: `<domain><devices><interface type='user'>
				<portForward proto='udp'><range start='2222' to='22'/></portForward>
				<portForward proto='tcp'><range start='8080' to='80'/><range start='2200' to='22'/></portForward>
			</interface></devices></domain>`,
			want: 2200,
		},
		{
			name:    "no forward",
			desc:    `<domain><devices><interface type='user'/></devices></domain>`,
			wantErr: true,
		},
		{
			name:    "malformed xml",
			desc:    `<domain><devices>`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseSSHPort(tt.desc)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseSSHPort() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("parseSSHPort() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestAddress(t *testing.T) {
	if got := Address(2222); got != "127.0.0.1:2222" {
		t.Fatalf("Address() = %q", got)
	}
}

func listenOnFreePort(t *testing.T) (net.Listener, int) {
	t.Helper()

	listener, err := net.Listen("tcp", Address(0))
	if err != nil {
		t.Fatalf("cannot listen: %v", err)
	}
	t.Cleanup(func() { _ = listener.Close() })

	_, portStr, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		t.Fatalf("cannot split address: %v", err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		t.Fatalf("cannot parse port: %v", err)
	}

	return listener, port
}

func TestEnsurePortFree(t *testing.T) {
	listener, port := listenOnFreePort(t)

	if err := EnsurePortFree(port); err == nil {
		t.Fatal("EnsurePortFree() succeeded on a bound port")
	}

	if err := listener.Close(); err != nil {
		t.Fatalf("cannot close listener: %v", err)
	}

	if err := EnsurePortFree(port); err != nil {
		t.Fatalf("EnsurePortFree() failed on a free port: %v", err)
	}
}

func TestIsSSHReachable(t *testing.T) {
	listener, port := listenOnFreePort(t)

	if !IsSSHReachable(port, time.Second) {
		t.Fatal("IsSSHReachable() = false for a listening port")
	}

	if err := listener.Close(); err != nil {
		t.Fatalf("cannot close listener: %v", err)
	}

	if IsSSHReachable(port, time.Second) {
		t.Fatal("IsSSHReachable() = true for a closed port")
	}
}
