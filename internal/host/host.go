package host

import (
	"bufio"
	"log"
	"net"
	"os/exec"
	"strings"

	"github.com/eikendev/hackenv/internal/paths"
)

func GetHostIPAddress(ifaceName string) string {
	ifaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Cannot retrieve host interfaces: %s\n", err)
	}

	for _, iface := range ifaces {
		if iface.Name != ifaceName {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}

			ip = ip.To16()
			if ip == nil {
				continue
			}

			return ip.String()
		}
	}

	log.Fatalln("Cannot retrieve host IP address")
	return "" // Does not actually return.
}

func GetHostKeyboardLayout() string {
	out, err := exec.Command(
		paths.GetCmdPathOrExit("setxkbmap"),
		"-query",
	).Output()
	if err != nil {
		log.Fatal(err)
	}

	var line string
	scanner := bufio.NewScanner(strings.NewReader(string(out)))

	for scanner.Scan() {
		line = scanner.Text()

		if strings.Contains(line, "layout") {
			break
		}
	}

	if line == "" {
		log.Fatalf("Unable to retrieve host's keyboard layout\n")
	}

	line = strings.TrimSpace(line)
	parts := strings.Split(line, " ")
	return parts[len(parts)-1]
}
