package upgrade

import (
	"fmt"

	"core/cli/fsdir"

	"github.com/blang/semver"
)

type Migration struct {
	Name     string
	Version  semver.Version
	Apply    func(fsdir.T) error
	Rollback func(fsdir.T) error
}

func (m *Migration) TryApply(f fsdir.T) error {
	if m != nil && m.Apply != nil {
		if err1 := m.Apply(f); err1 != nil {
			if m.Rollback != nil {
				if err2 := m.Rollback(f); err2 != nil {
					return fmt.Errorf("FAILED to roll back migration: %s, original error: %s", err2, err1)
				}
			}
			return fmt.Errorf("FAILED to apply migration: %s", err1)
		}
	}
	return nil
}
