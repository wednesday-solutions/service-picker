package utils

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func DirectoryName(dirName, stack, database string) string {
	switch stack {
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeHapiMySqlTemplate)
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeGraphqlMySqlTemplate)
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.NodeExpressMongoTemplate)
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangPgTemplate)
		} else if database == constants.MySQL {
			dirName = fmt.Sprintf("%s-%s", dirName, constants.GolangMySqlTemplate)
		}
	case constants.ReactJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactTemplate)
	case constants.NextJS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.NextTemplate)
	case constants.ReactNative:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactNativeTemplate)
	case constants.Android:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.AndroidTemplate)
	case constants.IOS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.IOSTemplate)
	case constants.Flutter:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.FlutterTemplate)
	}
	return dirName
}

func ExistingStacksDatabasesAndDirectories() ([]string, []string, []string) {
	var stacks, databases, dirNames []string
	var stack, database string
	directories, err := fileutils.ReadAllContents(fileutils.CurrentDirectory())
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

func FindStackAndDatabase(dirName string) (string, string) {
	var first, second, stack, database string
	splitDirName := strings.Split(dirName, "-")
	if len(splitDirName) > 2 {
		first = splitDirName[len(splitDirName)-1]
		second = splitDirName[len(splitDirName)-2]
		switch first {
		case "pg":
			database = constants.PostgreSQL
			if second == "hapi" {
				stack = constants.NodeHapiTemplate
			} else if second == "graphql" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if second == "golang" {
				stack = constants.GolangPostgreSQLTemplate
			}
		case "mysql":
			database = constants.MySQL
			if second == "hapi" {
				stack = constants.NodeHapiTemplate
			} else if second == "graphql" {
				stack = constants.NodeExpressGraphqlTemplate
			} else if second == "golang" {
				stack = constants.GolangMySQLTemplate
			}
		case "mongo":
			database = constants.MongoDB
			if second == "express" {
				stack = constants.NodeExpressTemplate
			}
		case "web":
			if second == "react" {
				stack = constants.ReactJS
			} else if second == "next" {
				stack = constants.NextJS
			}
		case "mobile":
			if second == "reactnative" {
				stack = constants.ReactNative
			} else if second == "android" {
				stack = constants.Android
			} else if second == "ios" {
				stack = constants.IOS
			} else if second == "flutter" {
				stack = constants.Flutter
			}
		}
	}
	return stack, database
}

func ExistingStackAndDatabase(dirName string) (string, string) {
	stack, database := FindStackAndDatabase(dirName)
	return stack, database
}

func FindService(dirName string) string {
	splitDirName := strings.Split(dirName, "-")
	if len(splitDirName) > 2 {
		suffix := splitDirName[len(splitDirName)-1]
		switch suffix {
		case "pg", "mysql":
			return constants.Backend
		case "web":
			return constants.Web
		case "mobile":
			return constants.Mobile
		}
	}
	return ""
}
