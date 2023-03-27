package helpers

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func ConvertTemplateDatabase(stack, database string, stackInfo map[string]interface{}) error {

	dbConfigFile := "/backend/config/db.js"
	isDatabaseSupported := true

	switch stack {
	case constants.NodeHapiTemplate:
		if database == constants.PostgreSQL {
			isDatabaseSupported = false
		}

	case constants.NodeExpressGraphqlTemplate:
		if database == constants.MySQL {
			isDatabaseSupported = false
		}
	default:
		return nil
	}

	if !isDatabaseSupported {
		// Add new dependencies to package.json
		err := UpdatePackageDotJson(stack)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = UpdateEnvFiles(stack)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		err = UpdateDBConfig(stack, dbConfigFile, stackInfo)
		errorhandler.CheckNilErr(err)

		// Convert queries
		err = ConvertQueries(stack)
		errorhandler.CheckNilErr(err)

		// Update docker-compose file
		err = UpdateDockerCompose(stack, stackInfo)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
