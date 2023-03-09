package fileutils

import (
	"os"
	"path/filepath"
)

// IsExists will check if the path exists or no.
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// CurrentDirectory will give the root directory.
func CurrentDirectory() string {
	path, err := filepath.Abs(".")
	if err != nil {
		return ""
	}
	return path
}

// MakeDirectory will make directory according to input.
func MakeDirectory(path string, dirName string) error {
	err := os.Mkdir(path+"/"+dirName, 0755)
	if err != nil {
		return err
	}
	return nil
}
