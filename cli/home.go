package cli

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"core/cli/fsdir"
)

// Home returns the fsdir of the home directory of this component.
func Home() fsdir.T {
	exe, err := os.Executable()

	if err != nil {
		log.Fatal(err)
	}

	fp, err := filepath.EvalSymlinks(exe)

	if err != nil {
		log.Fatal(err)
	}

	fm, err := fsdir.New(path.Dir(fp))

	if err != nil {
		log.Fatal(err)
	}

	return fm
}
