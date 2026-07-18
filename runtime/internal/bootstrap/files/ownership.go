package files

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
)

var chown = os.Chown

// OwnershipMatches reports whether a directory
// and all of its contents have  the requested ownership.
func OwnershipMatches(root string, uid int, gid int) (bool, error) {
	matches := true
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("unsupported file metadata for %q", path)
		}
		if uint64(stat.Uid) != uint64(uid) || uint64(stat.Gid) != uint64(gid) {
			matches = false
		}
		return nil
	})
	return matches, err
}

// ChownRecursive updates ownership for a directory and all of its contents
func ChownRecursive(root string, uid int, gid int) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		err = chown(path, uid, gid)
		if err != nil && os.IsNotExist(err) {
			// ignore entries removed while traversing directory tree
			return nil
		}
		return err
	})
}
