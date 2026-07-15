package files

import (
	"io/fs"
	"os"
	"path/filepath"
)

// ChownRecursive updates ownership for a directory and all of its contents
func ChownRecursive(root string, uid int, gid int) error {
	return filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chown(path, uid, gid)
	})
}
