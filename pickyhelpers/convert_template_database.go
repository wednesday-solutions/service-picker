package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

func (s StackDetails) ConvertTemplateDatabase() error {

	isDatabaseSupported := true

	switch s.Stack {
	case constants.NodeHapiTemplate:
		if s.Database == constants.PostgreSQL {
			isDatabaseSupported = false
		}

	case constants.NodeExpressGraphqlTemplate:
		if s.Database == constants.MySQL {
			isDatabaseSupported = false
		}
	default:
		return nil
	}

	if !isDatabaseSupported {
		// Add new dependencies to package.json
		err := UpdatePackageDotJson(s.Stack, s.DirName)
		errorhandler.CheckNilErr(err)

		// Update env files with respect to new database
		err = UpdateEnvFiles(s.Stack, s.DirName)
		errorhandler.CheckNilErr(err)

		// Convert DB Connection into MySQL.
		dbConfigFile := fmt.Sprintf("%s/%s/%s", utils.CurrentDirectory(), s.DirName, "config/db.js")
		err = UpdateDBConfig(s.Stack, dbConfigFile, s.StackInfo)
		errorhandler.CheckNilErr(err)

		// Convert queries
		err = ConvertQueries(s.Stack, s.DirName)
		errorhandler.CheckNilErr(err)

		// Update docker-compose file
		err = UpdateDockerCompose(s.Stack, s.DirName, s.StackInfo)
		errorhandler.CheckNilErr(err)
	} else {
		// Update DB_HOST in .env.docker file
		err := UpdateEnvDockerFileForDefaultDBInTemplate(s.Stack, s.DirName)
		errorhandler.CheckNilErr(err)
	}
	return nil
}
