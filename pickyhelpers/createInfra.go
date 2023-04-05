package pickyhelpers

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/hbs"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers/sources"
)

// IsInfraFilesExist checks the infra related files are exist of not.
func IsInfraFilesExist() bool {
	path := utils.CurrentDirectory()
	files := []string{
		constants.PackageDotJsonFile,
		constants.EnvFile,
		constants.SstConfigFile,
		constants.Stacks,
	}
	for _, file := range files {
		status, _ := utils.IsExists(filepath.Join(path, file))
		if !status {
			return false
		}
	}
	return true
}

// CreateInfraSetup creates package.json and .env files for infra setup.
func CreateInfraSetup() error {

	infraFiles := make(map[string]string)
	// package.json file
	infraFiles[constants.PackageDotJsonFile] = sources.PackageDotJsonSource()
	// .env file
	infraFiles[constants.EnvFile] = sources.EnvFileSource()

	var err error
	var path string
	for file, source := range infraFiles {
		path = fmt.Sprintf("%s/%s", utils.CurrentDirectory(), file)
		err = utils.WriteToFile(path, source)
		errorhandler.CheckNilErr(err)
	}
	return nil
}

// CreateInfraStacks creates the stack files of existing stacks in the stacks directory.
func CreateInfraStacks(service, stack, database, dirName, environment string) error {
	camelCaseDirName := strcase.ToCamel(dirName)
	var err error
	var stackFileName string
	path := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.Stacks)
	folderExist, _ := utils.IsExists(path)
	if !folderExist {
		err = utils.MakeDirectory(utils.CurrentDirectory(), constants.Stacks)
		errorhandler.CheckNilErr(err)
	}
	stackFileName = fmt.Sprintf("%s%s", camelCaseDirName, ".js")
	path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(), constants.Stacks, stackFileName)
	var source string

	switch service {
	case constants.Web:
		source = sources.WebStackSource(dirName, camelCaseDirName, environment)
	case constants.Mobile:
		// not implemented
	case constants.Backend:
		source = sources.BackendStackSource(database, dirName, environment)
	}
	err = utils.WriteToFile(path, source)
	return err
}

// CreateSstConfigFile creates sst.config.js file.
func CreateSstConfigFile(stackInfo map[string]interface{}, directories []string,
) error {
	sstConfigSource := sources.SstConfigSource()
	path := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.SstConfigFile)

	// SST config file for all selected existing stacks.
	camelCaseDirectories := utils.ConvertToCamelCase(directories)
	stackInfo[constants.ExistingDirectories] = camelCaseDirectories

	err := hbs.ParseAndWriteToFile(sstConfigSource, path, stackInfo)
	return err
}

// UpdateEnvByEnvironment updates the env file with respect to environment.
func UpdateEnvByEnvironment(dirName, environment string) error {
	path := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(),
		dirName,
		fmt.Sprintf("%s.%s", constants.EnvFile, environment),
	)
	err := utils.WriteToFile(path, sources.EnvSource(environment))
	errorhandler.CheckNilErr(err)
	return nil
}

// IsInfraStacksExist will return non existing stacks in the stacks directory.
// It will check the stack function is exists in the stacks directory.
func IsInfraStacksExist(stacksDirectories []string) []string {
	var path, camelCaseStack string
	var status bool
	var nonExistingStacks []string
	for _, stack := range stacksDirectories {
		camelCaseStack = fmt.Sprintf("%s%s", strcase.ToCamel(stack), ".js")
		path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(), constants.Stacks, camelCaseStack)
		status, _ = utils.IsExists(path)
		if !status {
			nonExistingStacks = append(nonExistingStacks, stack)
		}
	}
	return nonExistingStacks
}

// SstConfigExistStacks will give the stacks which is present in the sst.config.js
func SstConfigExistStacks() []string {
	file := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.SstConfigFile)
	status, _ := utils.IsExists(file)
	if status {
		// Reads the sst.config.js file and store the contents in input.
		input, err := ioutil.ReadFile(file)
		errorhandler.CheckNilErr(err)

		lines := strings.Split(string(input), "\n")
		// Pass the line which contain the stack details.
		// Eg: app.stack(ApiNodeHapiMysql).stack(ReReactWeb)
		// Will get the stacks which are present. Eg: [ApiNodeHapiMysql, FeReactWeb]
		configStackFiles := utils.FindConfigStacks(lines[len(lines)-4])

		// Will get stack directories. Eg: [api-node-hapi-mysql, fe-react-web]
		stacks := utils.FindStackDirectoriesByConfigStacks(configStackFiles)
		return stacks
	} else {
		return []string{}
	}
}
