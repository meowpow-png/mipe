package files

import "os"

// CreateDirectory creates a directory and any missing parent directories
func CreateDirectory(path string) error {
	return os.MkdirAll(path, 0o777)
}
