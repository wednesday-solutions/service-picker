package pickyhelpers

import (
	"fmt"
	"path/filepath"

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
		constants.ParseSstOutputs,
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

type Infra struct {
	Service          string
	Stack            string
	Database         string
	DirName          string
	CamelCaseDirName string
	Environment      string
	ForceCreate      bool
}

// CreateInfraStack creates the infra stack file of existing stack in the stacks directory.
func (i Infra) CreateInfraStack() error {
	var err error
	var stackFileName string
	path := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.Stacks)
	folderExist, _ := utils.IsExists(path)
	if !folderExist {
		err = utils.MakeDirectory(utils.CurrentDirectory(), constants.Stacks)
		errorhandler.CheckNilErr(err)
	}
	stackFileName = fmt.Sprintf("%s%s", i.CamelCaseDirName, ".js")
	path = fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(), constants.Stacks, stackFileName)
	response := true
	if !i.ForceCreate {
		stackFileExist, _ := utils.IsExists(path)
		if stackFileExist {
			return errorhandler.ErrExist
		}
	}
	var source string
	if response {
		switch i.Service {
		case constants.Web:
			source = sources.WebStackSource(i.DirName, i.CamelCaseDirName, i.Environment)
		case constants.Backend:
			source = sources.BackendStackSource(i.Database, i.DirName, i.Environment)
		default:
			err := utils.PrintErrorMessage("Selected stack is invalid")
			return err
		}
		err = utils.WriteToFile(path, source)
	}
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
	errorhandler.CheckNilErr(err)

	// create parseSstOutputs.js file
	parseSstOutputsSource := sources.ParseSstOutputsSource()
	path = fmt.Sprintf("%s/%s", utils.CurrentDirectory(), constants.ParseSstOutputs)

	err = utils.WriteToFile(path, parseSstOutputsSource)
	errorhandler.CheckNilErr(err)
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
func GetNonExistingInfraStacks(stacksDirectories []string) []string {
	var status bool
	var nonExistingStacks []string
	for _, stack := range stacksDirectories {
		status = IsInfraStackExist(stack)
		if !status {
			nonExistingStacks = append(nonExistingStacks, stack)
		}
	}
	return nonExistingStacks
}

func IsInfraStackExist(stackDirName string) bool {
	camelCaseStack := fmt.Sprintf("%s%s", strcase.ToCamel(stackDirName), ".js")
	path := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(), constants.Stacks, camelCaseStack)
	status, _ := utils.IsExists(path)
	return status
}
