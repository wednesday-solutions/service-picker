package utils

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
)

// CreateStackDirectory will create a a directory based on the user input, stack and the database selected.
func CreateStackDirectory(dirName, stack, database string) string {
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
	case constants.ReactGraphqlTS:
		dirName = fmt.Sprintf("%s-%s", dirName, constants.ReactGraphqlTemplate)
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

type StackDetails struct {
	Name      string
	Language  string
	Framework string
	Type      string
	Databases string
}

// GetStackDetails returns an array of StackDetails for showing details
// when user selects stacks prompt.
func GetStackDetails(service string) []StackDetails {
	var stacksDetails []StackDetails
	switch service {
	case constants.Backend:
		nodeHapi := StackDetails{
			Name:      constants.NodeHapiTemplate,
			Language:  "JavaScript",
			Framework: "Node JS & Hapi",
			Type:      "REST API",
			Databases: fmt.Sprintf("%s & %s", constants.PostgreSQL, constants.MySQL),
		}
		nodeGraphql := StackDetails{
			Name:      constants.NodeExpressGraphqlTemplate,
			Language:  "JavaScript",
			Framework: "Node JS & Express",
			Type:      "GraphQL API",
			Databases: fmt.Sprintf("%s & %s", constants.PostgreSQL, constants.MySQL),
		}
		nodeExpress := StackDetails{
			Name:      constants.NodeExpressTemplate,
			Language:  "JavaScript",
			Framework: "Node JS & Express",
			Type:      "REST API",
			Databases: constants.MongoDB,
		}
		golangGraphql := StackDetails{
			Name:      constants.GolangEchoTemplate,
			Language:  "Golang",
			Framework: "Echo",
			Type:      "GraphQL API",
			Databases: fmt.Sprintf("%s & %s", constants.PostgreSQL, constants.MySQL),
		}
		_, _, mobileStackExist := IsBackendWebAndMobileExist()
		if mobileStackExist {
			stacksDetails = []StackDetails{
				nodeHapi,
				nodeExpress,
			}
		} else {
			stacksDetails = []StackDetails{
				nodeHapi,
				nodeGraphql,
				nodeExpress,
				golangGraphql,
			}
		}
	case constants.Web:
		reactGraphqlTs := StackDetails{
			Name:      constants.ReactGraphqlTS,
			Language:  "TypeScript",
			Framework: "React",
			Type:      "GraphQL API",
		}
		reactJs := StackDetails{
			Name:      constants.ReactJS,
			Language:  "JavaScript",
			Framework: "React",
			Type:      "REST API",
		}
		nextJs := StackDetails{
			Name:      constants.NextJS,
			Language:  "JavaScript",
			Framework: "Next",
			Type:      "REST API",
		}
		stacksDetails = []StackDetails{reactJs, nextJs, reactGraphqlTs}

	case constants.Mobile:
		stacksDetails = []StackDetails{
			{
				Name:      constants.ReactNative,
				Language:  "JavaScript",
				Framework: "React Native",
				Type:      "REST API",
			},
			{
				Name:      constants.Android,
				Language:  "Kotlin",
				Framework: "-",
				Type:      "REST API",
			},
			{
				Name:      constants.IOS,
				Language:  "Swift",
				Framework: "-",
				Type:      "REST API",
			},
			{
				Name:      constants.Flutter,
				Language:  "Dart",
				Framework: "Flutter",
				Type:      "REST API",
			},
		}
	}
	return stacksDetails
}

// GetExistingInfraStacks fetch stack files inside the stacks directory.
func GetExistingInfraStacks() []string {
	path := fmt.Sprintf("%s/%s", CurrentDirectory(), constants.Stacks)
	status, _ := IsExists(path)
	if !status {
		return []string{}
	}
	files, err := ReadAllContents(path)
	errorhandler.CheckNilErr(err)
	return files
}

// GetSuffixOfStack returns suffix name for the given stack and database.
func GetSuffixOfStack(stack, database string) string {
	var suffix string
	switch stack {
	case constants.ReactJS:
		suffix = constants.ReactTemplate
	case constants.NextJS:
		suffix = constants.NextTemplate
	case constants.ReactGraphqlTS:
		suffix = constants.ReactGraphqlTemplate
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeHapiPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeHapiMySqlTemplate
		}
	case constants.NodeExpressGraphqlTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.NodeGraphqlPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.NodeGraphqlMySqlTemplate
		}
	case constants.NodeExpressTemplate:
		if database == constants.MongoDB {
			suffix = constants.NodeExpressMongoTemplate
		}
	case constants.GolangEchoTemplate:
		if database == constants.PostgreSQL {
			suffix = constants.GolangPgTemplate
		} else if database == constants.MySQL {
			suffix = constants.GolangMySqlTemplate
		}
	case constants.ReactNative:
		suffix = constants.ReactNativeTemplate
	case constants.Android:
		suffix = constants.AndroidTemplate
	case constants.IOS:
		suffix = constants.IOSTemplate
	case constants.Flutter:
		suffix = constants.FlutterTemplate
	}
	return suffix
}
