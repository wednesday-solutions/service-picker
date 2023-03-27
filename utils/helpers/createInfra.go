package helpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/hbs"
	"github.com/wednesday-solutions/picky/utils/helpers/sources"
)

func CreateInfra(stack, service string, stackInfo map[string]interface{}) error {

	infraFiles := make(map[string]string)
	path := fileutils.CurrentDirectory()

	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigJsFile,
		constants.WebStackJsFile,
		constants.BackendStackJsFile,
	}

	for _, file := range files {
		status, _ := fileutils.IsExists(path + "/" + file)
		if status {
			errorhandler.CheckNilErr(fmt.Errorf("%s file already exist", file))
		}
	}

	infraFiles[constants.PackageDotJsonFile] = sources.PackageDotJsonSource()

	infraFiles[constants.EnvFile] = sources.EnvFileSource()

	// SST config file for backend and web.
	infraFiles[constants.SstConfigJsFile] = sources.SstConfigJsSource()

	if stackInfo[constants.WebStatus].(bool) {
		// AWS config file for web.
		infraFiles[constants.WebStackJsFile] = sources.WebStackJsSource()
	}

	if stackInfo[constants.BackendStatus].(bool) {
		// AWS config file for backend.
		infraFiles[constants.BackendStackJsFile] = sources.BackendStackJsSource()
	}

	done := make(chan bool)
	go ProgressBar(30, "Generating", done)

	stacksFileExist, _ := fileutils.IsExists(fmt.Sprintf("%s/%s", path, constants.Stacks))
	for fileName, fileSource := range infraFiles {

		if fileName == constants.WebStackJsFile || fileName == constants.BackendStackJsFile {
			if !stacksFileExist {
				err := fileutils.MakeDirectory(path, "stacks")
				errorhandler.CheckNilErr(err)
				stacksFileExist = true
			}
			path = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.Stacks)

		} else {
			path = fileutils.CurrentDirectory()
		}
		var err error
		if fileName == constants.SstConfigJsFile {
			filePath := fmt.Sprintf("%s/%s", path, fileName)
			err = hbs.ParseAndWriteToFile(fileSource, filePath, stackInfo)
		} else {
			err = fileutils.TruncateAndWriteToFile(path, fileName, fileSource)
		}
		errorhandler.CheckNilErr(err)
	}

	// Update backend/.env.development file.
	path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
		constants.Backend,
		constants.EnvDevFile,
	)
	err := fileutils.WriteToFile(path, sources.EnvDevSource())
	errorhandler.CheckNilErr(err)

	<-done
	fmt.Printf("\nGenerating completed\n")

	return nil
}
