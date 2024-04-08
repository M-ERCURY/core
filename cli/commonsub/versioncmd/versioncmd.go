package versioncmd

import (
	"flag"
	"fmt"

	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/fsdir"
)

func Cmd(vstring string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("version", flag.ExitOnError),
		Desc:    "Show version and exit",
		Run:     func(_ fsdir.T) { fmt.Println(vstring) },
	}
}
