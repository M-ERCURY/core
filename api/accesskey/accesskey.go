package accesskey

import (
	"core/api/jsonb"
	"core/api/pof"
	"core/api/texturl"

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
