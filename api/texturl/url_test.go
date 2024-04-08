package texturl

import (
	"bytes"
	"testing"
)

func TestMarshalText(t *testing.T) {
	u := URLMustParse("https://mercury.com")
	r := []byte("https://mercury.com")

	m, err := u.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Compare(m, r) != 0 {
		t.Fatal("Marshalled text doesn't equal requested result")
	}
}

func TestUnmarshalText(t *testing.T) {
	var b []byte
	u := URLMustParse("https://mercury.com")

	if err := u.UnmarshalText(b); err != nil {
		t.Fatal(err)
	}
}

func TestURLMustParse(t *testing.T) {
	u := URLMustParse("https://mercury.com")
	if u.Scheme != "https" {
		t.Fatal("Did not parse URL properly")
	}
}
