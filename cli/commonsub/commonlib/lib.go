package commonlib

import (
	"github.com/M-ERCURY/core/cli/fsdir"

	"github.com/blang/semver"
)

// Context used by some commonsub commands for customization.
type Context struct {
	// Name of the binary using this command.
	BinName string
	// Target version to upgrade to.
	NewVersion semver.Version
	// Hook ran before command. Errors are logged and abort the execution.
	PreHook func(fsdir.T) error
	// Hook ran after command. Errors are logged and abort the execution.
	PostHook func(fsdir.T) error
}
