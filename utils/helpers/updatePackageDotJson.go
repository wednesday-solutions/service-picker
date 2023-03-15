package helpers

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdatePackageDotJson(stack, database string) error {

	var command string
	var dependencies []string
	var updateCommands []string

	switch stack {
	case constants.NODE_HAPI_TEMPLATE:
		// convert to postgres
		dependencies = []string{"pg", "pg-native"}

	case constants.NODE_EXPRESS_GRAPHQL_TEMPLATE:
		// convert to mysql
		dependencies = []string{"mysql2"}
	}

	err := exec.Command("yarn", "-v").Run()
	if err != nil {
		err = exec.Command("npm", "-v").Run()
		if err != nil {
			log.Fatal("Please install 'yarn' or 'npm' in your machine.")
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
	cmd.Dir = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), "backend")
	err = cmd.Run()
	errorhandler.CheckNilErr(err)

	return nil
}
