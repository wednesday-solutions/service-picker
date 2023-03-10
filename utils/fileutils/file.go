package fileutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

// IsExists will check if the path exists or no.
func IsExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	} else {
		return false, err
	}
}

// CurrentDirectory will give the root directory.
func CurrentDirectory() string {
	path, err := filepath.Abs(".")
	if err != nil {
		return ""
	} else {
		return path
	}
}

// MakeDirectory will make directory according to input.
func MakeDirectory(path string, dirName string) error {
	err := os.Mkdir(path+"/"+dirName, 0755)
	if err != nil {
		return err
	} else {
		return nil
	}
}

// MakeFile will create new file according to input path and file name.
func MakeFile(path, fileName string) error {
	_, err := os.Create(path + "/" + fileName)
	if err != nil {
		return err
	}
	return nil
}

func CreateFile(file string) error {
	_, err := os.Create(file)
	errorhandler.CheckNilErr(err)

	return nil
}

func RemoveFile(path string) error {
	err := os.Remove(path) // Remove a single file
	if err != nil {
		return err
	} else {
		return nil
	}
}

// TruncateAndWriteToFile will write input data into the file.
func TruncateAndWriteToFile(path, file, data string) error {
	// Opens file with read and write permissions.
	openFile, err := os.OpenFile(fmt.Sprintf("%s/%s", path, file), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer openFile.Close()

	_, err = openFile.WriteString(data)
	if err != nil {
		return err
	}

	err = openFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

// WriteToFile will write input data into the file.
func WriteToFile(file, data string) error {
	// Opens file with read and write permissions.
	openFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer openFile.Close()

	_, err = openFile.WriteString(data)
	if err != nil {
		return err
	}

	err = openFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

// AppendToFile will append given string to existing file.
func AppendToFile(path, file, data string) error {
	openFile, err := os.OpenFile(fmt.Sprintf("%s/%s", path, file), os.O_RDWR, 0644)
	errorhandler.CheckNilErr(err)

	defer openFile.Close()

	_, err = openFile.WriteString(data)
	errorhandler.CheckNilErr(err)

	err = openFile.Sync()
	errorhandler.CheckNilErr(err)

	return nil
}
