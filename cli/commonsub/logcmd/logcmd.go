package logcmd

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/M-ERCURY/core/cli"
	"github.com/M-ERCURY/core/cli/fsdir"
)

func Cmd(arg0 string) *cli.Subcmd {
	return &cli.Subcmd{
		FlagSet: flag.NewFlagSet("log", flag.ExitOnError),
		Desc:    fmt.Sprintf("Show %s logs", arg0),
		Run: func(fm fsdir.T) {
			logpath := fm.Path(arg0 + ".log")
			b, err := ioutil.ReadFile(logpath)
			if err != nil {
				log.Fatal(err)
			}
			os.Stdout.Write(b)
		},
	}
}
