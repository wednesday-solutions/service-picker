package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func PromptDockerCompose() {
	label := fmt.Sprintf("Do you want to create '%s' file for this project", constants.DockerComposeFile)
	response := PromptYesOrNoSelect(label)
	if response {
		GenerateDockerCompose()
	} else {
		PromptHome()
	}
}

func GenerateDockerCompose() {
	var stack, database, label string
	response := true
	status, _ := fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), constants.DockerComposeFile))
	if status {
		label = fmt.Sprintf("'%s' already exist, do you want to update it", constants.DockerComposeFile)
		response = PromptYesOrNoSelect(label)
	}
	if response {
		stacks, databases, _ := utils.ExistingStacksDatabasesAndDirectories()
		for i, db := range databases {
			if db != "" {
				database = db
				stack = stacks[i]
				break
			}
		}
		stackInfo := pickyhelpers.GetStackInfo(stack, database)
		err := pickyhelpers.CreateDockerComposeFile(stackInfo)
		errorhandler.CheckNilErr(err)
		fmt.Printf("\n%s\n", errorhandler.DoneMessage)
	}
	PromptHome()
}

func ShowCreateDockerCompose(databases []string) bool {
	var backendStatus, frontendStatus bool
	for _, db := range databases {
		if db == "" {
			frontendStatus = true
		} else {
			backendStatus = true
		}
		if backendStatus && frontendStatus {
			return true
		}
	}
	return false
}
