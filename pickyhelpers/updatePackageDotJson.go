package pickyhelpers

import (
	"fmt"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
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
	pkgManager = utils.GetPackageManagerOfUser()
	if pkgManager == constants.Yarn {
		updateCommands = []string{"add"}
		updateCommands = append(updateCommands, dependencies...)
	} else if pkgManager == constants.Npm {
		updateCommands = []string{"install", "--legacy-peer-deps", "--save"}
		updateCommands = append(updateCommands, dependencies...)
	}
	path := fmt.Sprintf("%s/%s", utils.CurrentDirectory(), dirName)
	err := utils.RunCommandWithoutLogs(path, pkgManager, updateCommands...)
	errorhandler.CheckNilErr(err)
	return nil
}
