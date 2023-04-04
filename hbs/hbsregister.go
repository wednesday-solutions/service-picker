package hbs

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
)

func DatabaseVolume(db string) string {
	if db == constants.PostgreSQL {
		return "postgresql/data"
	} else if db == constants.MySQL {
		return "mysql"
	} else {
		return ""
	}
}

func DBVersion(db string) string {
	if db == constants.PostgreSQL {
		return "postgres:15"
	} else if db == constants.MySQL {
		return "mysql:5.7"
	} else {
		return ""
	}
}

func PortConnection(stack string) string {
	switch stack {
	case constants.PostgreSQL:
		return "5432:5432"
	case constants.MySQL:
		return "3306:3306"
	case constants.MongoDB:
		return "27017:27017"
	case constants.Web, constants.Mobile:
		return "3000:3000"
	case constants.Backend:
		return "9000:9000"
	case constants.Redis:
		return "6379:6379"
	default:
		return ""
	}
}

func DBServiceName(stack, database string) string {
	switch stack {
	case constants.NodeExpressGraphqlTemplate, constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			return "db_postgres"
		} else if database == constants.MySQL {
			return "db_mysql"
		}
	case constants.GolangPostgreSQLTemplate, constants.GolangMySQLTemplate:
		return "db"
	}
	return "db"
}

func GlobalAddDependencies(database string) string {
	switch database {
	case constants.PostgreSQL, constants.MySQL:
		return "sequelize-cli@6.2.0"
	default:
		return ""
	}
}

func AddDependencies(database string) string {
	switch database {
	case constants.PostgreSQL:
		return "shelljs bull dotenv pg sequelize@6.6.5"
	case constants.MySQL:
		return "shelljs dotenv mysql2 sequelize@6.6.5"
	default:
		return ""
	}
}

func EnvFileBackend(database string) string {
	switch database {
	case constants.PostgreSQL, constants.MySQL:
		return "./backend/.env.docker"
	default:
		return ""
	}
}

func RunBuildEnvironment(stack string) string {
	switch stack {
	case constants.NodeExpressGraphqlTemplate:
		return "build:docker"
	case constants.NodeHapiTemplate:
		return "build:env"
	default:
		return ""
	}
}

func WaitForDBService(database string) string {
	var portNumber string
	if database == constants.PostgreSQL {
		portNumber = "5432"
	} else if database == constants.MySQL {
		portNumber = "3306"
	}
	return fmt.Sprintf(`  wait-for-db:
    image: atkrad/wait4x
    depends_on:
      - db
    command: tcp db:%s -t 30s -i 250ms`, portNumber)
}

func DependsOnFieldOfGo(stack string) string {
	output := `    depends_on:
      wait-for-db:
        condition: service_completed_successfully
`
	if stack == constants.GolangPostgreSQLTemplate || stack == constants.GolangMySQLTemplate {
		return output
	} else {
		return ""
	}
}

func CmdDockerfile(stack string) string {
	switch stack {
	case constants.ReactJS:
		return `["yarn", "start"]`
	case constants.NextJS:
		return `["yarn", "start:dev"]`
	default:
		return ""
	}
}

func EnvEnvironmentName() string {
	return "`.env.${process.env.ENVIRONMENT_NAME}`"
}

func DeployStacks(dirName string, directories []string) string {
	var source string
	if dirName == constants.All {
		for _, dir := range directories {
			source = fmt.Sprintf("%s.stack(%s)", source, dir)
		}
		source = fmt.Sprintf("app%s;", source)
	} else {
		source = fmt.Sprintf("app.stack(%s);", dirName)
	}
	return source
}

func SstImportStacks(dirName string, directories []string) string {
	var source string
	if dirName == constants.All {
		// import all existing stacks
		for _, dirName := range directories {
			source = fmt.Sprintf("%simport { %s } from %s./stacks/%s%s;\n", source, dirName, `"`, dirName, `"`)
		}
	} else {
		// import only one stack
		source = fmt.Sprintf("import { %s } from %s./stacks/%s%s;\n", dirName, `"`, dirName, `"`)
	}
	return source
}
