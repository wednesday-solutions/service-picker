package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
	"github.com/wednesday-solutions/picky/pickyhelpers"
)

func PromptDockerCompose() {
	var p PromptInput
	p.Label = fmt.Sprintf("Do you want to create '%s' file for this project", constants.DockerComposeFile)
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		GenerateDockerCompose()
	} else {
		PromptHome()
	}
}

func GenerateDockerCompose() {
	var p PromptInput
	p.Label = "Pick a service"
	p.GoBack = PromptDockerCompose
	var stack, database string
	response := true
	status, _ := utils.IsExists(filepath.Join(utils.CurrentDirectory(), constants.DockerComposeFile))
	if status {
		p.Label = fmt.Sprintf("'%s' already exist, do you want to update it", constants.DockerComposeFile)
		response = p.PromptYesOrNoSelect()
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
		stackInfo := pickyhelpers.GetStackInfo(stack, database, constants.Environment)
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
