package files

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyContents recursively copies all contents from source to destination
func CopyContents(source string, destination string) error {
	return copyChildren(source, destination)
}

func copyChildren(source string, destination string) error {
	entries, err := os.ReadDir(source)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(source, entry.Name())
		destinationPath := filepath.Join(destination, entry.Name())

		if err := copyPath(sourcePath, destinationPath); err != nil {
			return err
		}
	}
	return nil
}

func copyPath(source string, destination string) error {
	info, err := os.Stat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return copyDirectory(source, destination, info.Mode())
	}
	return copyFile(source, destination, info.Mode())
}

func copyDirectory(source string, destination string, mode fs.FileMode) error {
	if err := os.MkdirAll(destination, mode); err != nil {
		return err
	}
	return copyChildren(source, destination)
}

func copyFile(source string, destination string, mode fs.FileMode) (err error) {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := input.Close(); err == nil {
			err = closeErr
		}
	}()

	output, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := output.Close(); err == nil {
			err = closeErr
		}
	}()
	_, err = io.Copy(output, input)

	return err
}
