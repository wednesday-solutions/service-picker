package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func ConvertTemplateDatabase(stack, database, dirName string, stackInfo map[string]interface{}) error {

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
		err := UpdatePackageDotJson(stack, dirName)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = UpdateEnvFiles(stack, dirName)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		dbConfigFile := fmt.Sprintf("%s/%s/%s", fileutils.CurrentDirectory(), dirName, "config/db.js")
		err = UpdateDBConfig(stack, dbConfigFile, stackInfo)
		errorhandler.CheckNilErr(err)

		// Convert queries
		err = ConvertQueries(stack, dirName)
		errorhandler.CheckNilErr(err)

		// Update docker-compose file
		err = UpdateDockerCompose(stack, dirName, stackInfo)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
