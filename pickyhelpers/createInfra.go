package pickyhelpers

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func IsInfraFilesExist() bool {
	path := fileutils.CurrentDirectory()
	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigJsFile,
		constants.Stacks,
	}
	for _, file := range files {
		status, _ := fileutils.IsExists(filepath.Join(path, file))
		if status {
			return true
		}
	}
	return false
}

func CreateInfraSetup(stackInfo map[string]interface{}) error {

	infraFiles := make(map[string]string)
	// package.json file
	infraFiles[constants.PackageDotJsonFile] = sources.PackageDotJsonSource()
	// .env file
	infraFiles[constants.EnvFile] = sources.EnvFileSource()

	var err error
	var path string
	for file, source := range infraFiles {
		path = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), file)
		if file == constants.SstConfigJsFile {
			err = hbs.ParseAndWriteToFile(source, path, stackInfo)
		} else {
			err = fileutils.WriteToFile(path, source)
		}
		errorhandler.CheckNilErr(err)
	}
	return nil
}

func CreateInfraStacks(service, stack, database, dirName string, stackInfo map[string]interface{}) error {
	var err error
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.Stacks)
	folderExist, _ := fileutils.IsExists(path)
	if !folderExist {
		err = fileutils.MakeDirectory(fileutils.CurrentDirectory(), constants.Stacks)
		errorhandler.CheckNilErr(err)
	}
	dirName = fmt.Sprintf("%s%s", dirName, ".js")
	path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(), constants.Stacks, dirName)
	var source string
	status, _ := fileutils.IsExists(path)
	if !status {
		switch service {
		case constants.Web:
			source = sources.WebStackJsSource()
		case constants.Mobile:
			// not implemented
		case constants.Backend:
			source = sources.BackendStackJsSource()
		}
		err = fileutils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)

		// SST config file for backend and web.
		sstConfigSource := sources.SstConfigJsSource()
		path = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.SstConfigJsFile)
		err = hbs.ParseAndWriteToFile(sstConfigSource, path, stackInfo)
		errorhandler.CheckNilErr(err)
	} else {
		return errorhandler.ErrExist
	}
	return nil
}

func UpdateEnvDevelopment(dirName string) error {
	path := fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
		dirName,
		constants.EnvDevFile,
	)
	err := fileutils.WriteToFile(path, sources.EnvDevSource())
	errorhandler.CheckNilErr(err)
	return nil
}
