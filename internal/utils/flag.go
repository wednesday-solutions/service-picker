package utils

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/internal/constants"
)

func UseService() string {
	usageString := fmt.Sprintf(`choose a service
  %d. %s
  %d. %s
  %d. %s
`,
		1, constants.Web,
		2, constants.Mobile,
		3, constants.Backend)
	return usageString
}

func AllWebStacksString() string {
	usage := fmt.Sprintf(`
	Web stacks:
	 %d. %s          -> %s
	 %d. %s           -> %s
	 %d. %s  -> %s
`,
		1, constants.ReactJS, constants.ReactjsLower,
		2, constants.NextJS, constants.NextjsLower,
		3, constants.ReactGraphqlTS, constants.ReactGraphqlLower,
	)
	return usage
}

func AllMobileStacksString() string {
	usage := fmt.Sprintf(`
	Mobile stacks:
	 %d. %s -> %s
	 %d. %s      -> %s
	 %d. %s          -> %s
	 %d. %s      -> %s
`,
		1, constants.ReactNative, constants.ReactNativeLower,
		2, constants.Android, constants.AndroidLower,
		3, constants.IOS, constants.IOSLower,
		4, constants.Flutter, constants.FlutterLower,
	)
	return usage
}
func AllBackendStacksString() string {
	usage := fmt.Sprintf(`
	Backend stacks:
	 %d. %s       -> %s
	 %d. %s -> %s
	 %d. %s    -> %s
	 %d. %s  -> %s
`,
		1, constants.NodeHapiTemplate, constants.NodeHapi,
		2, constants.NodeExpressGraphqlTemplate, constants.NodeGraphql,
		3, constants.NodeExpressTemplate, constants.NodeExpress,
		4, constants.GolangEchoTemplate, constants.Golang,
	)
	return usage
}

func AllStacksString() string {
	usage := fmt.Sprintf(`%s%s%s`,
		AllWebStacksString(),
		AllMobileStacksString(),
		AllBackendStacksString(),
	)
	return usage
}

func UseStack() string {
	usageString := fmt.Sprintf("choose a stack (select the second name)\n%s", AllStacksString())
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
	return `provide a directory prefix name (suffix will be added)
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
	_, _, directories := GetExistingStacksDatabasesAndDirectories()
	var usageString string
	if len(directories) == 0 {
		usageString = "Stacks not exist. Existing stacks see here.\n"
		return usageString
	}
	usageString = "existing stacks are\n"
	for idx, dir := range directories {
		usageString = fmt.Sprintf("%s %d. %s\n", usageString, idx+1, dir)
	}
	return usageString
}

func GetExistingStacks() []string {
	_, _, directories := GetExistingStacksDatabasesAndDirectories()
	return directories
}

func UseCloudProvider() string {
	usageString := fmt.Sprintf(`choose a cloud provider
 %d. %s
`, 1, constants.AWS)
	return usageString
}

func UseEnvironment() string {
	usageString := fmt.Sprintf(`choose an environment
 %d. %s
 %d. %s
 %d. %s
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
	if env == constants.Development || env == constants.Develop || env == constants.Dev {
		return constants.Development
	}
	return env
}

func UseDockerCompose() string {
	return "create Docker Compose file for the stacks which are present in the root directory."
}

func UseCI() string {
	return "create CI file for the stacks which are present in the root directory."
}

func UseCD() string {
	return "create CD file for the stacks which are present in the root directory."
}

func UsePlatform() string {
	usage := fmt.Sprintf(`choose a platform
  %d. %s
`,
		1, constants.Github)
	return usage
}

func ConvertStacksIntoString(stacks []string) string {
	var response string
	for idx, stack := range stacks {
		response = fmt.Sprintf("%s\n  %d. %s", response, idx+1, stack)
	}
	return response
}
