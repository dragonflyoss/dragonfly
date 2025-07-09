//go:build windows
// +build windows

package storage

import (
	"os"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func hardlink(log *logger.SugaredLoggerOnWith, dst, src string) error {
	// Try to create a hard link if the destination does not exist
	err := error(nil)
	if _, statErr := os.Stat(dst); os.IsNotExist(statErr) {
		err = os.Link(src, dst)
		if err != nil {
			log.Errorf("hardlink from %q to %q error: %s", src, dst, err)
			return err
		}
		log.Infof("hardlink from %q to %q success", src, dst)
		return nil
	}

	if err != nil {
		log.Errorf("stat %q error: %s", src, err)
		return err
	}

	// On Windows, inode check is skipped. Return error if file exists.
	return err
}
