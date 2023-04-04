package prompt

import (
	"github.com/wednesday-solutions/picky/pickyhelpers"
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
		err := pickyhelpers.UpdatePackageDotJson(stack)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = pickyhelpers.UpdateEnvFiles(stack)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		err = pickyhelpers.UpdateDBConfig(stack, dbConfigFile, stackInfo)
		errorhandler.CheckNilErr(err)

		// Convert queries
		err = pickyhelpers.ConvertQueries(stack)
		errorhandler.CheckNilErr(err)

		// Update docker-compose file
		err = pickyhelpers.UpdateDockerCompose(stack, stackInfo)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
