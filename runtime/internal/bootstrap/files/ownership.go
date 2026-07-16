package files

import (
	"io/fs"
	"os"
	"path/filepath"
)

var chown = os.Chown

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
