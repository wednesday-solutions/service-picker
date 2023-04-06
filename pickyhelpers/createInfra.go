package pickyhelpers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func IsInfraFilesExist() bool {
	path := fileutils.CurrentDirectory()
	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigFile,
		constants.Stacks,
	}
	for _, file := range files {
		status, _ := fileutils.IsExists(filepath.Join(path, file))
		if !status {
			return false
		}
	}
	return true
}

func CreateInfraSetup() error {

	infraFiles := make(map[string]string)
	// package.json file
	infraFiles[constants.PackageDotJsonFile] = sources.PackageDotJsonSource()
	// .env file
	infraFiles[constants.EnvFile] = sources.EnvFileSource()

	var err error
	var path string
	for file, source := range infraFiles {
		path = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), file)
		err = fileutils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)
	}
	return nil
}

func CreateInfraStacks(service, stack, database, dirName, environment string) error {
	camelCaseDirName := strcase.ToCamel(dirName)
	var err error
	var stackFileName string
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.Stacks)
	folderExist, _ := fileutils.IsExists(path)
	if !folderExist {
		err = fileutils.MakeDirectory(fileutils.CurrentDirectory(), constants.Stacks)
		errorhandler.CheckNilErr(err)
	}
	stackFileName = fmt.Sprintf("%s%s", camelCaseDirName, ".js")
	path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(), constants.Stacks, stackFileName)
	var source string

	switch service {
	case constants.Web:
		source = sources.WebStackSource(camelCaseDirName, environment)
	case constants.Mobile:
		// not implemented
	case constants.Backend:
		source = sources.BackendStackSource(database, dirName, environment)
	}
	err = fileutils.WriteToFile(path, source)
	return err
}

func CreateSstConfigFile(stackInfo map[string]interface{}, directories []string,
) error {
	sstConfigSource := sources.SstConfigSource()
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.SstConfigFile)

	// SST config file for all selected existing stacks.
	camelCaseDirectories := utils.ToCamelCase(directories)
	stackInfo[constants.ExistingDirectories] = camelCaseDirectories

	err := hbs.ParseAndWriteToFile(sstConfigSource, path, stackInfo)
	return err
}

func UpdateEnvDevelopment(dirName, environment string) error {
	path := fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(),
		dirName,
		constants.EnvDevFile,
	)
	err := fileutils.WriteToFile(path, sources.EnvDevSource(environment))
	errorhandler.CheckNilErr(err)
	return nil
}

func IsInfraStacksExist(services []string) []string {
	var path, camelCaseService string
	var status bool
	var nonExistingStacks []string
	for _, service := range services {
		camelCaseService = fmt.Sprintf("%s%s", strcase.ToCamel(service), ".js")
		path = fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(), constants.Stacks, camelCaseService)
		status, _ = fileutils.IsExists(path)
		if !status {
			nonExistingStacks = append(nonExistingStacks, service)
		}
	}
	return nonExistingStacks
}

func SstConfigExistStacks() []string {
	input, err := ioutil.ReadFile("sst.config.js")
	errorhandler.CheckNilErr(err)

	lines := strings.Split(string(input), "\n")
	configStackFiles := utils.FindConfigStacks(lines[len(lines)-4])

	stacks := utils.FindStacksByConfigStacks(configStackFiles)
	return stacks
}
