package utils

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/internal/constants"
)

func UseService() string {
	usageString := fmt.Sprintf(`Choose a service
  %d. %s
  %d. %s
  %d. %s
`,
		1, constants.Web,
		2, constants.Mobile,
		3, constants.Backend)
	return usageString
}

func UseStack() string {
	usageString := fmt.Sprintf(`Choose a stack (select the second name)

 Web stacks:
  %d. %s          -> %s
  %d. %s           -> %s
  %d. %s  -> %s

 Mobile stacks:
  %d. %s -> %s
  %d. %s      -> %s
  %d. %s          -> %s
  %d. %s      -> %s

 Backend stacks:
  %d. %s       -> %s
  %d. %s -> %s
  %d. %s    -> %s
  %d. %s  -> %s
`,
		1, constants.ReactJS, constants.ReactjsLower,
		2, constants.NextJS, constants.NextjsLower,
		3, constants.ReactGraphqlTS, constants.ReactGraphqlLower,

		1, constants.ReactNative, constants.ReactNativeLower,
		2, constants.Android, constants.AndroidLower,
		3, constants.IOS, constants.IOSLower,
		4, constants.Flutter, constants.FlutterLower,

		1, constants.NodeHapiTemplate, constants.NodeHapi,
		2, constants.NodeExpressGraphqlTemplate, constants.NodeGraphql,
		3, constants.NodeExpressTemplate, constants.NodeExpress,
		4, constants.GolangEchoTemplate, constants.Golang)
	return usageString
}

func UseDatabase() string {
	usageString := fmt.Sprintf(`Choose a database
  %d. %s
  %d. %s
  %d. %s
`,
		1, constants.Postgresql,
		2, constants.Mysql,
		3, constants.Mongodb,
	)
	return usageString
}

func UseDirectory() string {
	return `Provide a directory name (suffix will be added.)
  Eg: example-react-js-web | example-node-hapi-pg`
}

func GetStackByFlags(stack string) string {
	switch stack {
	case constants.ReactjsLower, constants.ReactJS:
		return constants.ReactJS
	case constants.NextjsLower, constants.NextJS:
		return constants.NextJS
	case constants.ReactGraphqlLower, constants.ReactGraphqlTS:
		return constants.ReactGraphqlTS

	case constants.ReactNativeLower, constants.ReactNative:
		return constants.ReactNative
	case constants.AndroidLower, constants.Android:
		return constants.Android
	case constants.IOSLower, constants.IOS:
		return constants.IOS
	case constants.Flutter, constants.FlutterLower:
		return constants.FlutterLower

	case constants.NodeHapi, constants.NodeHapiTemplate:
		return constants.NodeHapiTemplate
	case constants.NodeGraphql, constants.NodeExpressGraphqlTemplate:
		return constants.NodeExpressGraphqlTemplate
	case constants.NodeExpress, constants.NodeExpressTemplate:
		return constants.NodeExpressTemplate
	case constants.Golang, constants.GolangEchoTemplate:
		return constants.GolangEchoTemplate

	default:
		return ""
	}
}

func GetDatabase(db string) string {
	db = strings.ToLower(db)
	if db == constants.Postgresql || db == constants.Postgres {
		return constants.PostgreSQL
	} else if db == constants.Mysql {
		return constants.MySQL
	} else if db == constants.Mongodb {
		return constants.MongoDB
	} else {
		return db
	}
}

func UseInfraStacks() string {
	_, _, directories := ExistingStacksDatabasesAndDirectories()
	var usageString string
	if len(directories) == 0 {
		usageString = "Stacks not exist. Existing stacks see here.\n"
		return usageString
	}
	usageString = "Existing stacks are\n"
	for idx, dir := range directories {
		usageString = fmt.Sprintf("%s %d. %s\n", usageString, idx+1, dir)
	}
	return usageString
}

func ExistingStacks() []string {
	_, _, directories := ExistingStacksDatabasesAndDirectories()
	return directories
}

func UseCloudProvider() string {
	usageString := fmt.Sprintf(`Choose a cloud provider
 %d. %s
aws is the default cloud provider.
`, 1, constants.AWS)
	return usageString
}

func UseEnvironment() string {
	usageString := fmt.Sprintf(`Choose an environment
 %d. %s
 %d. %s
 %d. %s
development is the default environment.
`, 1, constants.Development,
		2, constants.QA,
		3, constants.Production)
	return usageString
}

func GetCloudProvider(cp string) string {
	cp = strings.ToLower(cp)
	if cp == "aws" {
		return constants.AWS
	}
	return cp
}

func GetEnvironmentValue(env string) string {
	env = strings.ToLower(env)
	if env == "development" || env == "develop" || env == "dev" {
		return constants.Development
	}
	return env
}
