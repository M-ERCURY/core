package accesskey

import (
	"github.com/M-ERCURY/core/api/jsonb"
	"github.com/M-ERCURY/core/api/pof"
	"github.com/M-ERCURY/core/api/texturl"

	"github.com/blang/semver"
)

type T struct {
	Version  *semver.Version `json:"version"`
	Contract *Contract       `json:"contract,omitempty"`
	Pofs     []*pof.T        `json:"pofs,omitempty"`
}

type Contract struct {
	Endpoint  *texturl.URL `json:"endpoint,omitempty"`
	PublicKey jsonb.PK     `json:"public_key,omitempty"`
}
