package images

import (
	"bytes"
	_ "embed"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
)

// Kali Linux Archive Automatic Signing Key (2025) <devel@kali.org>
// Obtained from https://archive.kali.org/archive-key.asc
//
//go:embed keys/kali.asc
var kaliSigningKey []byte

// Parrot Project (2024-2026) <team@parrotsec.org>
// Obtained from https://deb.parrot.sh/parrot/misc/parrotsec.gpg
//
//go:embed keys/parrot.asc
var parrotSigningKey []byte

// SigningKey describes an embedded vendor key used to verify checksum files.
type SigningKey struct {
	Image       string
	Fingerprint string
	URL         string
	key         []byte
}

var signingKeys = []SigningKey{
	{Image: "kali", Fingerprint: "827C8569F2518CC677FECA1AED65462EC8D5E4C5", URL: "https://archive.kali.org/archive-key.asc", key: kaliSigningKey},
	{Image: "parrot", Fingerprint: "B711822346552E4D92DA02DF7A8286AF0E81EE4A", URL: "https://deb.parrot.sh/parrot/misc/parrotsec.gpg", key: parrotSigningKey},
}

// SigningKeys returns the embedded vendor keys.
func SigningKeys() []SigningKey {
	return signingKeys
}

// Expiry returns when the embedded key expires, or the zero time if it never expires.
func (k SigningKey) Expiry() (time.Time, error) {
	keyring, err := readSigningKey(k.key)
	if err != nil {
		return time.Time{}, err
	}

	primary := keyring[0].PrimaryKey
	identity := keyring[0].PrimaryIdentity()
	if identity == nil || identity.SelfSignature == nil || identity.SelfSignature.KeyLifetimeSecs == nil {
		return time.Time{}, nil
	}

	return primary.CreationTime.Add(time.Duration(*identity.SelfSignature.KeyLifetimeSecs) * time.Second), nil
}

// ActualFingerprint returns the fingerprint of the embedded key material.
func (k SigningKey) ActualFingerprint() (string, error) {
	keyring, err := readSigningKey(k.key)
	if err != nil {
		return "", err
	}

	return strings.ToUpper(hex.EncodeToString(keyring[0].PrimaryKey.Fingerprint)), nil
}

func readSigningKey(armored []byte) (openpgp.EntityList, error) {
	keyring, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(armored))
	if err != nil {
		return nil, fmt.Errorf("cannot read signing key: %w", err)
	}
	if len(keyring) == 0 {
		return nil, errors.New("signing key is empty")
	}

	return keyring, nil
}

// verifyDetached checks a detached signature, armored or binary, over data and returns the data.
func verifyDetached(armoredKey, data, signature []byte) ([]byte, error) {
	keyring, err := readSigningKey(armoredKey)
	if err != nil {
		return nil, err
	}

	check := openpgp.CheckDetachedSignature
	if bytes.HasPrefix(bytes.TrimSpace(signature), []byte("-----BEGIN PGP")) {
		check = openpgp.CheckArmoredDetachedSignature
	}

	if _, err := check(keyring, bytes.NewReader(data), bytes.NewReader(signature), nil); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return data, nil
}

// verifyClearsigned checks a clearsigned message and returns the signed plaintext.
func verifyClearsigned(armoredKey, message []byte) ([]byte, error) {
	keyring, err := readSigningKey(armoredKey)
	if err != nil {
		return nil, err
	}

	block, _ := clearsign.Decode(message)
	if block == nil {
		return nil, errors.New("message is not clearsigned")
	}

	if _, err := block.VerifySignature(keyring, nil); err != nil {
		return nil, fmt.Errorf("signature verification failed: %w", err)
	}

	return block.Plaintext, nil
}
