package images

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	rawLibvirt "libvirt.org/go/libvirt"
)

type genericVersionComparer struct{}

func getGenericVersionComparer() *genericVersionComparer {
	return &genericVersionComparer{}
}

func (vc genericVersionComparer) Lt(a, b string) bool {
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")

	if len(aParts) == 0 || len(bParts) == 0 || len(aParts) != len(bParts) {
		slog.Error("Cannot compare versions with different parts", "a", a, "b", b)
		panic(fmt.Sprintf("Cannot compare versions %s and %s", a, b))
	}

	for i := range aParts {
		aPart, err := strconv.Atoi(aParts[i])
		if err != nil {
			slog.Error("Cannot convert version part to number", "value", aParts[i], "err", err)
			panic(fmt.Sprintf("Cannot convert version part %s to number", aParts[i]))
		}

		bPart, err := strconv.Atoi(bParts[i])
		if err != nil {
			slog.Error("Cannot convert version part to number", "value", bParts[i], "err", err)
			panic(fmt.Sprintf("Cannot convert version part %s to number", bParts[i]))
		}

		if aPart < bPart {
			return true
		}
		if aPart > bPart {
			return false
		}
	}

	return false
}

func (vc genericVersionComparer) Eq(a, b string) bool {
	return a == b
}

func (vc genericVersionComparer) Gt(a, b string) bool {
	return !vc.Lt(a, b) && !vc.Eq(a, b)
}

func genericBootInitializer(dom *rawLibvirt.Domain) error {
	time.Sleep(1 * time.Second)
	return sendKeys(dom, []uint{KEY_ENTER})
}

func switchToTTY(dom *rawLibvirt.Domain) error {
	if err := sendKeys(dom, []uint{KEY_LEFTCTRL, KEY_LEFTALT, KEY_F1}); err != nil {
		slog.Error("Failed to send key sequence to switch to TTY", "err", err)
		return fmt.Errorf("failed to switch console to TTY: %w", err)
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func switchFromTTY(dom *rawLibvirt.Domain) error {
	if err := sendKeys(dom, []uint{KEY_LEFTCTRL, KEY_LEFTALT, KEY_F7}); err != nil {
		slog.Error("Failed to send key sequence to switch from TTY", "err", err)
		return fmt.Errorf("failed to switch console from TTY: %w", err)
	}
	time.Sleep(500 * time.Millisecond)
	return nil
}

func enablePasswordSSH(dom *rawLibvirt.Domain) error {
	// sudo sed -i '/.assword.uthentication/s/no/yes/' /etc/ssh/sshd*<Tab>

	keys := []uint{
		KEY_S, KEY_U, KEY_D, KEY_O, KEY_SPACE, KEY_S, KEY_E, KEY_D, KEY_SPACE,
		KEY_MINUS, KEY_I, KEY_SPACE, KEY_APOSTROPHE, KEY_SLASH, KEY_DOT, KEY_A,
		KEY_S, KEY_S, KEY_W, KEY_O, KEY_R, KEY_D, KEY_DOT, KEY_U, KEY_T, KEY_H,
		KEY_E, KEY_N, KEY_T, KEY_I, KEY_C, KEY_A, KEY_T, KEY_I, KEY_O, KEY_N,
		KEY_SLASH, KEY_S, KEY_SLASH, KEY_N, KEY_O, KEY_SLASH, KEY_Y, KEY_E,
		KEY_S, KEY_SLASH, KEY_APOSTROPHE, KEY_SPACE, KEY_SLASH, KEY_E, KEY_T,
		KEY_C, KEY_SLASH, KEY_S, KEY_S, KEY_H, KEY_SLASH, KEY_S, KEY_S, KEY_H,
		KEY_D, KEY_TAB, KEY_ENTER,
	}

	for _, key := range keys {
		if err := sendKeys(dom, []uint{key}); err != nil {
			slog.Error("Failed to send key while enabling password SSH", "err", err, "key", key)
			return fmt.Errorf("failed to enable password SSH: %w", err)
		}
	}

	time.Sleep(500 * time.Millisecond)
	return nil
}

func systemdRestartSSH(dom *rawLibvirt.Domain) error {
	// sudo systemctl restart ssh

	keys := []uint{
		KEY_S, KEY_U, KEY_D, KEY_O, KEY_SPACE, KEY_S, KEY_Y, KEY_S, KEY_T,
		KEY_E, KEY_M, KEY_C, KEY_T, KEY_L, KEY_SPACE, KEY_R, KEY_E, KEY_S,
		KEY_T, KEY_A, KEY_R, KEY_T, KEY_SPACE, KEY_S, KEY_S, KEY_H, KEY_ENTER,
	}

	for _, key := range keys {
		if err := sendKeys(dom, []uint{key}); err != nil {
			slog.Error("Failed to send key while restarting SSH via systemd", "err", err, "key", key)
			return fmt.Errorf("failed to restart SSH via systemd: %w", err)
		}
	}

	time.Sleep(1500 * time.Millisecond)
	return nil
}
