package images

import (
	"time"

	rawLibvirt "libvirt.org/libvirt-go"
)

func genericBootInitializer(dom *rawLibvirt.Domain) {
	time.Sleep(1 * time.Second)
	sendKeys(dom, []uint{KEY_ENTER})
}

func switchToTTY(dom *rawLibvirt.Domain) {
	// Includes two dummy keys. See [0].
	// [0] https://gitlab.com/libvirt/libvirt-go/-/issues/10
	sendKeys(dom, []uint{KEY_LEFTCTRL, KEY_LEFTALT, KEY_F1, KEY_RESERVED, KEY_RESERVED})

	time.Sleep(500 * time.Millisecond)
}

func switchFromTTY(dom *rawLibvirt.Domain) {
	// Includes two dummy keys as above.
	sendKeys(dom, []uint{KEY_LEFTCTRL, KEY_LEFTALT, KEY_F7, KEY_RESERVED, KEY_RESERVED})

	time.Sleep(500 * time.Millisecond)
}

func enablePasswordSSH(dom *rawLibvirt.Domain) {
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
		sendKeys(dom, []uint{key})
	}

	time.Sleep(500 * time.Millisecond)
}

func systemdRestartSSH(dom *rawLibvirt.Domain) {
	// sudo systemctl restart ssh

	keys := []uint{
		KEY_S, KEY_U, KEY_D, KEY_O, KEY_SPACE, KEY_S, KEY_Y, KEY_S, KEY_T,
		KEY_E, KEY_M, KEY_C, KEY_T, KEY_L, KEY_SPACE, KEY_R, KEY_E, KEY_S,
		KEY_T, KEY_A, KEY_R, KEY_T, KEY_SPACE, KEY_S, KEY_S, KEY_H, KEY_ENTER,
	}

	for _, key := range keys {
		sendKeys(dom, []uint{key})
	}

	time.Sleep(1500 * time.Millisecond)
}
