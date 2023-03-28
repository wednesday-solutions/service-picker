package prompt

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/helpers"
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
		err := helpers.UpdatePackageDotJson(stack)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = helpers.UpdateEnvFiles(stack)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		err = helpers.UpdateDBConfig(stack, dbConfigFile, stackInfo)
		errorhandler.CheckNilErr(err)

		// Convert queries
		err = helpers.ConvertQueries(stack)
		errorhandler.CheckNilErr(err)

		// Update docker-compose file
		err = helpers.UpdateDockerCompose(stack, stackInfo)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
