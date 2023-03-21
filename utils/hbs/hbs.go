package hbs

import (
	"github.com/aymerick/raymond"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func init() {
	raymond.RegisterHelper("databaseVolume", DatabaseVolume)
	raymond.RegisterHelper("dbVersion", DBVersion)
	raymond.RegisterHelper("portConnection", PortConnection)
	raymond.RegisterHelper("dbServiceName", DBServiceName)
	raymond.RegisterHelper("globalAddDependencies", GlobalAddDependencies)
	raymond.RegisterHelper("addDependencies", AddDependencies)
	raymond.RegisterHelper("envFileBackend", EnvFileBackend)
	raymond.RegisterHelper("runBuildEnvironment", RunBuildEnvironment)
}

func ParseAndWriteToFile(source, filePath string, stackData map[string]interface{}) error {

	ctx := map[string]interface{}{
		constants.Frontend: constants.Frontend,
		constants.Web:      constants.Web,
		constants.Mobile:   constants.Mobile,
		constants.Backend:  constants.Backend,
		constants.Redis:    constants.Redis,
		constants.Postgres: constants.Postgres,
		constants.Mysql:    constants.Mysql,
		"stack":            stackData["stack"].(string),
		"database":         stackData["database"].(string),
		"projectName":      stackData["projectName"].(string),
		"webStatus":        stackData["webStatus"].(bool),
		"mobileStatus":     stackData["mobileStatus"].(bool),
		"backendStatus":    stackData["backendStatus"].(bool),
	}
	// Parse the source string into template
	tpl, err := raymond.Parse(source)
	errorhandler.CheckNilErr(err)

	// Execute the template into string
	executedTemplate, err := tpl.Exec(ctx)
	errorhandler.CheckNilErr(err)

	err = fileutils.WriteToFile(filePath, executedTemplate)
	errorhandler.CheckNilErr(err)

	return nil
}

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
	case "redis":
		return "6379:6379"
	default:
		return ""
	}
}

func DBServiceName(stack string) string {
	switch stack {
	case constants.PostgreSQL:
		return "db_postgres"
	case constants.MySQL:
		return "db_mysql"
	default:
		return "db"
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
