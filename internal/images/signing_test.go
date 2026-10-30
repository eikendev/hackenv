package images

import (
	"bytes"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
)

func newTestKey(t *testing.T) (*openpgp.Entity, []byte) {
	t.Helper()

	entity, err := openpgp.NewEntity("hackenv test", "", "test@example.invalid", nil)
	if err != nil {
		t.Fatalf("cannot create key: %v", err)
	}

	var armored bytes.Buffer
	w, err := armor.Encode(&armored, openpgp.PublicKeyType, nil)
	if err != nil {
		t.Fatalf("cannot start armor: %v", err)
	}
	if err := entity.Serialize(w); err != nil {
		t.Fatalf("cannot serialize key: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("cannot finish armor: %v", err)
	}

	return entity, armored.Bytes()
}

func detachSign(t *testing.T, entity *openpgp.Entity, data []byte) []byte {
	t.Helper()

	var signature bytes.Buffer
	if err := openpgp.DetachSign(&signature, entity, bytes.NewReader(data), nil); err != nil {
		t.Fatalf("cannot sign: %v", err)
	}

	return signature.Bytes()
}

func armoredDetachSign(t *testing.T, entity *openpgp.Entity, data []byte) []byte {
	t.Helper()

	var signature bytes.Buffer
	if err := openpgp.ArmoredDetachSign(&signature, entity, bytes.NewReader(data), nil); err != nil {
		t.Fatalf("cannot sign: %v", err)
	}

	return signature.Bytes()
}

func clearSign(t *testing.T, entity *openpgp.Entity, data []byte) []byte {
	t.Helper()

	var message bytes.Buffer
	w, err := clearsign.Encode(&message, entity.PrivateKey, nil)
	if err != nil {
		t.Fatalf("cannot start clearsign: %v", err)
	}
	if _, err := w.Write(data); err != nil {
		t.Fatalf("cannot write clearsign: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("cannot finish clearsign: %v", err)
	}

	return message.Bytes()
}

func TestVerifyDetached(t *testing.T) {
	entity, key := newTestKey(t)
	_, otherKey := newTestKey(t)
	data := []byte("abc123  kali-linux-2026.3-live-amd64.iso\n")
	signature := detachSign(t, entity, data)

	if got, err := verifyDetached(key, data, signature); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("verifyDetached() = %q, %v", got, err)
	}

	if got, err := verifyDetached(key, data, armoredDetachSign(t, entity, data)); err != nil || !bytes.Equal(got, data) {
		t.Fatalf("verifyDetached() with armored signature = %q, %v", got, err)
	}

	if _, err := verifyDetached(key, []byte("tampered"), signature); err == nil {
		t.Fatal("verifyDetached() accepted tampered data")
	}

	if _, err := verifyDetached(otherKey, data, signature); err == nil {
		t.Fatal("verifyDetached() accepted a signature from another key")
	}

	if _, err := verifyDetached([]byte("not a key"), data, signature); err == nil {
		t.Fatal("verifyDetached() accepted an invalid key")
	}
}

func TestVerifyClearsigned(t *testing.T) {
	entity, key := newTestKey(t)
	_, otherKey := newTestKey(t)
	data := []byte("sha256\nabc123  Parrot-security-7.4_amd64.iso\n")
	message := clearSign(t, entity, data)

	got, err := verifyClearsigned(key, message)
	if err != nil {
		t.Fatalf("verifyClearsigned() error = %v", err)
	}
	if !bytes.Contains(got, []byte("Parrot-security-7.4_amd64.iso")) {
		t.Fatalf("verifyClearsigned() plaintext = %q", got)
	}

	if _, err := verifyClearsigned(otherKey, message); err == nil {
		t.Fatal("verifyClearsigned() accepted a signature from another key")
	}

	tampered := bytes.Replace(message, []byte("abc123"), []byte("abc124"), 1)
	if _, err := verifyClearsigned(key, tampered); err == nil {
		t.Fatal("verifyClearsigned() accepted tampered data")
	}

	if _, err := verifyClearsigned(key, data); err == nil {
		t.Fatal("verifyClearsigned() accepted an unsigned message")
	}
}

func TestEmbeddedKeysMatchDocumentedFingerprints(t *testing.T) {
	for _, key := range SigningKeys() {
		actual, err := key.ActualFingerprint()
		if err != nil {
			t.Fatalf("embedded %s key unreadable: %v", key.Image, err)
		}
		if actual != key.Fingerprint {
			t.Fatalf("embedded %s key fingerprint = %s, documented %s", key.Image, actual, key.Fingerprint)
		}

		expiry, err := key.Expiry()
		if err != nil {
			t.Fatalf("embedded %s key expiry: %v", key.Image, err)
		}
		if expiry.IsZero() {
			t.Fatalf("embedded %s key has no expiry", key.Image)
		}
	}
}
