package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// SplitStackDirectoryName returns user-input, stack-suffix and last-suffix of the given stack directory name.
func SplitStackDirectoryName(dirName string) (string, string, string, string) {
	var userInput, langSuffix, stackSuffix, lastSuffix string
	var splitDirName []string
	var isBackendStack, isWebStack, isMobileStack bool
	splitDirName = strings.Split(dirName, "-")
	if len(splitDirName) > constants.Two {
		lastSuffix = splitDirName[len(splitDirName)-constants.One]
		stackSuffix = splitDirName[len(splitDirName)-constants.Two]
		langSuffix = splitDirName[len(splitDirName)-constants.Three]
		if lastSuffix == constants.Pg || lastSuffix == constants.Mysql || lastSuffix == constants.Mongo {
			isBackendStack = true
		} else if lastSuffix == constants.Web {
			isWebStack = true
		} else if lastSuffix == constants.Mobile {
			isMobileStack = true
		}
		var suffixSize int
		if isBackendStack {
			suffixSize = constants.BackendSuffixSize
		} else if isWebStack {
			suffixSize = constants.WebSuffixSize
		} else if isMobileStack {
			suffixSize = constants.MobileSuffixSize
		}
		userInput = splitDirName[constants.Zero]
		for _, split := range splitDirName[constants.One : len(splitDirName)-suffixSize] {
			userInput = fmt.Sprintf("%s_%s", userInput, split)
		}
	}
	return userInput, langSuffix, stackSuffix, lastSuffix
}

// CovertToCamelCase return camel cased array of string of the given array of string.
func ConvertToCamelCase(slice []string) []string {
	camelSlice := []string{}
	for _, str := range slice {
		camelSlice = append(camelSlice, strcase.ToCamel(str))
	}
	return camelSlice
}

// RunCommandWithLogs runs the given command with logs.
func RunCommandWithLogs(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if path != "" {
		cmd.Dir = path
	}
	err := cmd.Run()
	fmt.Printf("\n")
	return err
}

// RunCommandWithLogs runs the given command without logs.
func RunCommandWithoutLogs(path string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if path != "" {
		cmd.Dir = path
	}
	err := cmd.Run()
	return err
}

// GetPackageManagerOfUser checks whether yarn or npm is installed in the user's machine.
// If both are not installed, then the system will throw error.
func GetPackageManagerOfUser() string {
	var pkgManager string
	err := RunCommandWithoutLogs("", constants.Yarn, "-v")
	if err != nil {
		err = RunCommandWithoutLogs("", constants.Npm, "-v")
		if err != nil {
			// Throw error either yarn or npm not installed
			errorhandler.CheckNilErr(fmt.Errorf("Please install 'yarn' or 'npm' in your machine.\n"))
		} else {
			pkgManager = constants.Npm
		}
	} else {
		pkgManager = constants.Yarn
	}
	return pkgManager
}

// GetEnvironment return short environment name for the given environment.
func GetEnvironment(environment string) string {
	switch environment {
	case constants.Development:
		return constants.Dev
	case constants.QA:
		return constants.QA
	case constants.Production:
		return constants.Prod
	default:
		return environment
	}
}

// GetProjectName returns projectName
func GetProjectName() string {
	currentDir := CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	return projectName
}

func CreateGithubWorkflowDir() {
	currentDir := CurrentDirectory()
	workflowsPath := fmt.Sprintf("%s/%s", currentDir,
		constants.GithubWorkflowsDir,
	)
	workflowStatus, _ := IsExists(workflowsPath)
	if !workflowStatus {
		githubFolderPath := fmt.Sprintf("%s/%s", currentDir, ".github")
		githubStatus, _ := IsExists(githubFolderPath)
		if !githubStatus {
			err := CreateDirectory(githubFolderPath)
			errorhandler.CheckNilErr(err)
		}
		err := CreateDirectory(workflowsPath)
		errorhandler.CheckNilErr(err)
	}
}
