package relay

import (
	"bytes"
	"crypto/ed25519"
	"net"
	"testing"
	"time"

	"github.com/M-ERCURY/core/api/servicekey"
	"github.com/M-ERCURY/core/api/sharetoken"
	"github.com/M-ERCURY/core/api/signer"
	"github.com/M-ERCURY/core/api/texturl"
	"github.com/M-ERCURY/core/mrnet"
	"github.com/M-ERCURY/core/mrnet/transport"
)

func TestSMRelay(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Error(err)
	}
	tt := transport.New(transport.Options{
		TLSVerify: false,
		Timeout:   time.Second * 5,
	})
	rl := New(tt, Options{
		MaxTime:       time.Second * 0,
		BufSize:       2048,
		AllowLoopback: true,
	})

	s := signer.New(priv)
	sk := servicekey.New(priv)
	var (
		now    = time.Now()
		sopen  = now.Add(1 * time.Minute)
		sclose = sopen.Add(1 * time.Minute)
	)

	sk.Contract.SettlementOpen = sopen.Unix()
	sk.Contract.SettlementClose = sclose.Unix()
	sk.Contract.Sign(s)
	st, err := sharetoken.New(sk, pub)
	if err != nil {
		t.Fatal(err)
	}
	init := &mrnet.Init{
		Command:  "CONNECT",
		Protocol: "tcp",
		Remote:   texturl.URLMustParse("target://localhost:8888"),
		Token:    st,
		Version:  &mrnet.PROTO_VERSION,
	}
	p0 := []byte{'h', 'e', 'l', 'l', 'o', '!', '\r', '\n'}
	c1, c2 := net.Pipe()
	// emulate relay
	go rl.ServeTLS(c2)
	// emulate target
	l, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		t.Fatal(err)
	}
	// emulate client
	err = init.WriteTo(c1)
	if err != nil {
		t.Fatal(err)
	}
	con, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	n, err := c1.Write(p0)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p0) {
		t.Fatal("partial write")
	}
	p1 := make([]byte, 32)
	n, err = con.Read(p1)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(p0) {
		t.Fatal("partial read")
	}
	if !bytes.Equal(p0, p1[:n]) {
		t.Fatal("target received corrupted message:", p0, p1[:n])
	}
	_, err = con.Write(p1[:n])
	if err != nil {
		t.Fatal(err)
	}
	con.Close()
	c1 = &mrnet.FragReadConn{Conn: c1}
	p2 := make([]byte, 32)
	n, err = c1.Read(p2)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(p0, p2[:n]) {
		t.Fatal("mercury-relay received corrupted message", p0, p2[:n])
	}
}
