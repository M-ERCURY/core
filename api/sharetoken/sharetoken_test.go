package sharetoken

import (
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/M-ERCURY/core/api/jsonb"
	"github.com/M-ERCURY/core/api/servicekey"
	"github.com/M-ERCURY/core/api/signer"
)

func TestVerify(t *testing.T) {
	pk, sk, err := ed25519.GenerateKey(nil)

	if err != nil {
		t.Fatal(err)
	}

	skey := servicekey.New(sk)

	// mock contract sig
	skey.Contract = &servicekey.Contract{
		PublicKey:       jsonb.PK(pk),
		SettlementOpen:  9999999999,
		SettlementClose: 99999999999,
	}

	skey.Contract.Sign(signer.New(sk))

	sharetoken, err := New(skey, pk)

	if err != nil {
		t.Fatal(err)
	}

	err = sharetoken.Verify()

	if err != nil {
		t.Fatal(err)
	}

	if sharetoken.IsExpiredAt(time.Now().Unix()) {
		t.Fatal("sharetoken is expired")
	}
}
