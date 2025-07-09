//go:build unix
// +build unix

package storage

import (
	"fmt"
	"os"
	"syscall"

	logger "d7y.io/dragonfly/v2/internal/dflog"
)

func hardlink(log *logger.SugaredLoggerOnWith, dst, src string) error {
	dstStat, err := os.Stat(dst)
	if os.IsNotExist(err) {
		// hard link
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

	// target already exists, check inode
	srcStat, err := os.Stat(src)
	if err != nil {
		log.Errorf("stat %q error: %s", src, err)
		return err
	}

	dstSysStat, ok := dstStat.Sys().(*syscall.Stat_t)
	if !ok {
		log.Errorf("can not get inode for %q", dst)
		return err
	}

	srcSysStat, ok := srcStat.Sys().(*syscall.Stat_t)
	if !ok {
		log.Errorf("can not get inode for %q", src)
		return err
	}

	if dstSysStat.Dev == srcSysStat.Dev && dstSysStat.Ino == srcSysStat.Ino {
		log.Debugf("target inode match underlay data inode, skip hard link")
		return nil
	}

	err = fmt.Errorf("target file %q exists, with different inode with underlay data %q", dst, src)
	return err
}
