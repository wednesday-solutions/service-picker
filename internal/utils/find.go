package utils

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// FindStackAndDatabase return stack and database of given directory name.
func FindStackAndDatabase(dirName string) (string, string) {
	var stack, database string
	_, langSuffix, stackSuffix, lastSuffix := SplitStackDirectoryName(dirName)

	switch lastSuffix {
	case constants.Pg:
		database = constants.PostgreSQL
		if stackSuffix == "hapi" {
			stack = constants.NodeHapiTemplate
		} else if stackSuffix == constants.Graphql {
			if langSuffix == "node" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if langSuffix == "golang" {
				stack = constants.GolangEchoTemplate
			}
		}
	case constants.Mysql:
		database = constants.MySQL
		if stackSuffix == "hapi" {
			stack = constants.NodeHapiTemplate
		} else if stackSuffix == constants.Graphql {
			if langSuffix == "node" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if langSuffix == "golang" {
				stack = constants.GolangEchoTemplate
			}
		}
	case constants.Web:
		if stackSuffix == "js" {
			if langSuffix == "react" {
				stack = constants.ReactJS
			} else if langSuffix == "next" {
				stack = constants.NextJS
			}
		} else if stackSuffix == constants.Graphql {
			if langSuffix == "react" {
				stack = constants.ReactGraphqlTS
			}
		}
	case constants.Mobile:
		if stackSuffix == "reactnative" {
			stack = constants.ReactNative
		} else if stackSuffix == "android" {
			stack = constants.Android
		} else if stackSuffix == "ios" {
			stack = constants.IOS
		} else if stackSuffix == "flutter" {
			stack = constants.Flutter
		}
	}
	return stack, database
}

// ExistingStacksDatabasesAndDirectories return existing stacks details.
// Stack details contain stack name, database, and the directory name.
func ExistingStacksDatabasesAndDirectories() ([]string, []string, []string) {
	var stacks, databases, dirNames []string
	var stack, database string
	directories, err := ReadAllContents(CurrentDirectory())
	errorhandler.CheckNilErr(err)

	for _, dirName := range directories {
		stack, database = FindStackAndDatabase(dirName)
		if stack != "" {
			stacks = append(stacks, stack)
			databases = append(databases, database)
			dirNames = append(dirNames, dirName)
		}
	}
	return stacks, databases, dirNames
}

// FindUserInputStackName return user-input of the given stack directory name.
func FindUserInputStackName(dirName string) string {
	userInput, _, _, _ := SplitStackDirectoryName(dirName)
	return userInput
}

// ExistingStackAndDatabase return stack and database of the given stack directory name.
func ExistingStackAndDatabase(dirName string) (string, string) {
	stack, database := FindStackAndDatabase(dirName)
	return stack, database
}

// FindService return service of the given stack directory name.
func FindService(dirName string) string {
	_, _, _, lastSuffix := SplitStackDirectoryName(dirName)
	switch lastSuffix {
	case constants.Pg, constants.Mysql, constants.Mongo:
		return constants.Backend
	default:
		return lastSuffix
	}
}

// FindStackDirectoriesByConfigStacks will return an array of stack directories present
// in the root directory. Eg: [api-node-hapi-mysql, fe-react-web]
func FindStackDirectoriesByConfigStacks(configStacks []string) []string {
	var stacks []string

	_, _, directories := ExistingStacksDatabasesAndDirectories()
	camelCaseDirectories := ConvertToCamelCase(directories)

	for _, configStack := range configStacks {
		for idx, camelCaseDirName := range camelCaseDirectories {
			camelCaseDirName = fmt.Sprintf("%s%s", camelCaseDirName, ".js")
			if configStack == camelCaseDirName {
				stacks = append(stacks, directories[idx])
			}
		}
	}
	return stacks
}

// IsBackendWebAndMobileExist checks backend, web and mobile stacks is exist or not.
func IsBackendWebAndMobileExist() (bool, bool, bool) {
	var service string
	backendExist, webExist, mobileExist := false, false, false
	_, _, stackDirectories := ExistingStacksDatabasesAndDirectories()
	for _, stackDir := range stackDirectories {
		service = FindService(stackDir)
		if service == constants.Web {
			webExist = true
		} else if service == constants.Backend {
			backendExist = true
		} else if service == constants.Mobile {
			mobileExist = true
		}
	}
	return backendExist, webExist, mobileExist
}

func FindExistingGraphqlBackendAndWebStacks() ([]string, []string) {
	var backendGraphqlStacks, webGraphqlStacks []string
	var stack, service string
	_, _, stackDirectories := ExistingStacksDatabasesAndDirectories()
	for _, stackDir := range stackDirectories {
		service = FindService(stackDir)
		if service == constants.Backend {
			stack, _ = FindStackAndDatabase(stackDir)
			if stack == constants.NodeExpressGraphqlTemplate ||
				stack == constants.GolangEchoTemplate {
				backendGraphqlStacks = append(backendGraphqlStacks, stackDir)
			}
		} else if service == constants.Web {
			stack, _ = FindStackAndDatabase(stackDir)
			if stack == constants.ReactGraphqlTS {
				webGraphqlStacks = append(webGraphqlStacks, stackDir)
			}
		}
	}
	return backendGraphqlStacks, webGraphqlStacks
}
