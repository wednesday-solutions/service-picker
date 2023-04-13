package hbs

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func DatabaseVolumeConnection(db string) string {
	if db == constants.PostgreSQL {
		return fmt.Sprintf("-db-volume:/var/lib/%s", constants.PostgresqlData)
	} else if db == constants.MySQL {
		return fmt.Sprintf("-db-volume:/var/lib/%s", constants.Mysql)
	} else {
		return db
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
	var portConnectionStr string
	switch stack {
	case constants.PostgreSQL:
		portConnectionStr = fmt.Sprintf("%d:5432", utils.PostgresPortNumber)
		utils.PostgresPortNumber++
		return portConnectionStr
	case constants.MySQL:
		portConnectionStr = fmt.Sprintf("%d:3306", utils.MysqlPortNumber)
		utils.MysqlPortNumber++
		return portConnectionStr
	case constants.MongoDB:
		return "27017:27017"
	case constants.Web, constants.Mobile:
		portConnectionStr = fmt.Sprintf("%d:3000", utils.WebPortNumber)
		utils.WebPortNumber++
		return portConnectionStr
	case constants.Backend:
		portConnectionStr = fmt.Sprintf("%d:9000", utils.BackendPortNumber)
		utils.BackendPortNumber++
		return portConnectionStr
	case constants.Redis:
		return "6379:6379"
	default:
		return portConnectionStr
	}
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

func DeployStacks(stackFiles []string) string {
	var deployStackSource string
	for _, stackFile := range stackFiles {
		// will append all the selected stack files to deploy in sst.config.js
		deployStackSource = fmt.Sprintf("%s.stack(%s)", deployStackSource, stackFile)
	}
	deployStackSource = fmt.Sprintf("app%s;", deployStackSource)
	return deployStackSource
}

func SstImportStacks(stackFiles []string) string {
	var importStackSource string
	// import all existing stacks in sst.config.js
	for _, stackFile := range stackFiles {
		importStackSource = fmt.Sprintf("%simport { %s } from %s./stacks/%s%s;\n",
			importStackSource, stackFile, `"`, stackFile, `"`)
	}
	return importStackSource
}
