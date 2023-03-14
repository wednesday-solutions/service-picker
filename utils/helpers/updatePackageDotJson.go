package helpers

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func UpdatePackageDotJson(stack, database string) error {

	fmt.Println("Loading... (please wait, this will take some time)")

	var command string
	updateCommands := []string{}
	err := exec.Command("yarn", "-v").Run()
	if err != nil {
		log.Fatal("Command ", err)
		err = exec.Command("npm", "-v").Run()
		if err != nil {
			log.Fatal("Please install 'yarn' or 'npm' in your machine.")
		} else {
			command = "npm"
			updateCommands = append(updateCommands, "install", "--legacy-peer-deps", "--save", "pg", "pg-native")
		}
	} else {
		command = "yarn"
		updateCommands = append(updateCommands, "add", "pg", "pg-native")
	}

	cmd := exec.Command(command, updateCommands...)
	cmd.Dir = fmt.Sprintf("%s/%s", fileutils.CurrentDirectory(), "backend")
	err = cmd.Run()
	errorhandler.CheckNilErr(err)

	return nil
}
