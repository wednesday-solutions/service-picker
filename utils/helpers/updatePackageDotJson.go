package helpers

import (
	"fmt"
	"os/exec"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdatePackageDotJson(stack string) error {

	var command string
	var dependencies []string
	var updateCommands []string

	switch stack {
	case constants.NodeHapiTemplate:
		// convert to postgres
		dependencies = []string{"pg", "pg-native"}

	case constants.NodeExpressGraphqlTemplate:
		// convert to mysql
		dependencies = []string{"mysql2"}
	}

	err := exec.Command("yarn", "-v").Run()
	if err != nil {
		err = exec.Command("npm", "-v").Run()
		if err != nil {
			errorhandler.CheckNilErr(fmt.Errorf("Please install 'yarn' or 'npm' in your machine."))
		} else {
			command = "npm"
			updateCommands = []string{"install", "--legacy-peer-deps", "--save"}
			updateCommands = append(updateCommands, dependencies...)
		}
	} else {
		command = "yarn"
		updateCommands = []string{"add"}
		updateCommands = append(updateCommands, dependencies...)
	}

	cmd := exec.Command(command, updateCommands...)
	cmd.Dir = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), constants.Backend)
	err = cmd.Run()
	errorhandler.CheckNilErr(err)

	return nil
}
