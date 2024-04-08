package pof

import (
	"crypto/ed25519"
	"testing"
	"time"

	"core/api/signer"
)

func TestPof(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	p, err := New(signer.New(priv), "test", 100)

	if err != nil {
		t.Fatal(err)
	}

	if !p.IsExpiredAt(time.Now().Unix() + 100) {
		t.Fatal("incorrect expiry")
	}
}
