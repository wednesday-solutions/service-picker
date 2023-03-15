package helpers

import (
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func ConvertTemplateDatabase(stack, database, projectName string) error {

	dbConfigFile := "/backend/config/db.js"
	isDatabaseSupported := true

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:
		if database == constants.POSTGRES {
			isDatabaseSupported = false
		}

	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:
		if database == constants.MYSQL {
			isDatabaseSupported = false
		}
	default:
		return nil
	}

	if !isDatabaseSupported {
		// Add new dependencies to package.json
		err := UpdatePackageDotJson(stack, database)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = UpdateEnvFiles(stack, database, projectName)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		err = UpdateDBConfig(stack, dbConfigFile, database, projectName)
		errorhandler.CheckNilErr(err)

		// Convert mysql queries to postgres queries
		err = ConvertMysqlToPostgres(stack, database, projectName)
		errorhandler.CheckNilErr(err)
	}

	return nil
}
