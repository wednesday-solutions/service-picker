package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdatePackageDotJson(stack, dirName string) error {
	var pkgManager string
	var dependencies []string
	var updateCommands []string

	switch stack {
	case constants.NodeHapiTemplate:
		// convert to postgres
		dependencies = []string{constants.Pg, "pg-native"}

	case constants.NodeExpressGraphqlTemplate:
		// convert to mysql
		dependencies = []string{"mysql2"}
	}
	pkgManager = utils.IsYarnOrNpmInstalled()
	if pkgManager == constants.Yarn {
		updateCommands = []string{"add"}
		updateCommands = append(updateCommands, dependencies...)
	} else if pkgManager == constants.Npm {
		updateCommands = []string{"install", "--legacy-peer-deps", "--save"}
		updateCommands = append(updateCommands, dependencies...)
	}
	path := fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), dirName)
	err := utils.RunCommandWithoutLogs(path, pkgManager, updateCommands...)
	errorhandler.CheckNilErr(err)
	return nil
}
