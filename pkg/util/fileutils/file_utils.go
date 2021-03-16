/*
 *     Copyright 2020 The Dragonfly Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package fileutils provides utilities supplementing the standard about file packages.
package fileutils

import (
	"d7y.io/dragonfly/v2/pkg/util/fileutils/fsize"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"

	logger "d7y.io/dragonfly/v2/pkg/dflog"
	"github.com/pkg/errors"
)

// MkdirAll creates a directory named path with 0755 perm.
func MkdirAll(dir string) error {
	return os.MkdirAll(dir, 0755)
}

// DeleteFile deletes a regular file not a directory.
func DeleteFile(path string) error {
	if PathExist(path) {
		if IsDir(path) {
			return errors.Errorf("delete file %s: not a regular file", path)
		}

		return os.Remove(path)
	}

	return nil
}

// OpenFile opens a file. If the parent directory of the file isn't exist,
// it will create the directory.
func OpenFile(path string, flag int, perm os.FileMode) (*os.File, error) {
	if !PathExist(path) && (flag&syscall.O_CREAT != 0) {
		if err := MkdirAll(filepath.Dir(path)); err != nil {
			return nil, errors.Wrapf(err, "failed to open file %s", path)
		}
	}

	return os.OpenFile(path, flag, perm)
}

func Open(path string) (*os.File, error) {
	return OpenFile(path, syscall.O_RDONLY, 0)
}

// Link creates a hard link pointing to oldname named newname for a file.
func Link(oldname string, newname string) error {
	if PathExist(newname) {
		if IsDir(newname) {
			return errors.Errorf("link %s to %s: src already exists and is a directory", newname, oldname)
		}

		if err := DeleteFile(newname); err != nil {
			return errors.Wrapf(err, "failed to link %s to %s", newname, oldname)
		}
	} else if err := MkdirAll(filepath.Dir(newname)); err != nil {
		return errors.Wrapf(err, "failed to link %s to %s", newname, oldname)
	}

	return os.Link(oldname, newname)
}

func IsSymbolicLink(name string) bool {
	f, e := os.Lstat(name)
	if e != nil {
		return false
	}
	return f.Mode() & os.ModeSymlink != 0
}

// SymbolicLink creates newname as a symbolic link to oldname.
func SymbolicLink(oldname string, newname string) error {
	if PathExist(newname) && IsSymbolicLink(newname) && PathExist(oldname){
		dstPath, err := os.Readlink(newname)
		if err != nil {
			return err
		}
		if dstPath == oldname {
			return nil
		}
	}
	if PathExist(newname) {
		if err := os.Remove(newname); err != nil {
			return fmt.Errorf("failed to symlink %s to %s when deleting target file: %v", newname, oldname, err)
		}
	}

	if err := MkdirAll(filepath.Dir(newname)); err != nil {
		return errors.Wrapf(err, "failed to symlink %s to %s", newname, oldname)
	}
	return os.Symlink(oldname, newname)
}

// PathExist reports whether the path is exist.
// Any error, from os.Lstat, will return false.
func PathExist(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}

func IsDir(path string) bool {
	f, err := stat(path)
	if err != nil {
		return false
	}

	return f.IsDir()
}

func IsRegular(path string) bool {
	f, err := stat(path)
	if err != nil {
		return false
	}

	return f.Mode().IsRegular()
}

func IsEmptyDir(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	if _, err = f.Readdirnames(1); err == io.EOF {
		return true, nil
	}

	return false, err
}

func stat(path string) (os.FileInfo, error) {
	f, err := os.Stat(path)

	if err != nil {
		logger.Warnf("stat file %s: %v", path, err)
	}

	return f, err
}

// GetFreeSpace gets the free disk space of the path.
func GetFreeSpace(path string) (fsize.Size, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, err
	}

	return fsize.Size(fs.Bavail * uint64(fs.Bsize)), nil
}

// GetTotalSpace
func GetTotalSpace(path string) (fsize.Size, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, err
	}

	return fsize.Size(fs.Blocks * uint64(fs.Bsize)), nil
}

func GetTotalAndFreeSpace(path string) (fsize.Size, fsize.Size, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, 0, err
	}
	total := fsize.Size(fs.Blocks * uint64(fs.Bsize))
	free := fsize.Size(fs.Bavail * uint64(fs.Bsize))
	return total, free, nil
}

// GetUsedSpace
func GetUsedSpace(path string) (fsize.Size, error) {
	fs := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &fs); err != nil {
		return 0, err
	}
	return fsize.Size((fs.Blocks - fs.Bavail) * uint64(fs.Bsize)), nil
}

// MoveFile moves the file src to dst.
func MoveFile(src string, dst string) error {
	if !IsRegular(src) {
		return fmt.Errorf("failed to move %s to %s: src is not a regular file", src, dst)
	}
	if PathExist(dst) && !IsDir(dst) {
		if err := DeleteFile(dst); err != nil {
			return fmt.Errorf("failed to move %s to %s when deleting dst file: %v", src, dst, err)
		}
	}
	return os.Rename(src, dst)
}