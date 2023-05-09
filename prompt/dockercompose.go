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
	p.Label = "Choose an option"
	p.GoBack = PromptHome
	p.Items = []string{constants.CreateDockerCompose, constants.RunDockerCompose}
	response, _ := p.PromptSelect()
	if response == constants.CreateDockerCompose {
		err := GenerateDockerCompose()
		errorhandler.CheckNilErr(err)
		PromptRunDockerCompose()
	} else if response == constants.RunDockerCompose {
		err := RunDockerCompose()
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

// PromptCreateDockerCompose is a prompt function for create docker-compose.
func PromptCreateDockerCompose() {
	var p PromptInput
	p.Label = fmt.Sprintf("Do you want to create '%s' file for this project", constants.DockerComposeFile)
	p.GoBack = PromptHome
	response := p.PromptYesOrNoSelect()
	if response {
		err := GenerateDockerCompose()
		errorhandler.CheckNilErr(err)
	} else {
		PromptDockerCompose()
	}
	PromptRunDockerCompose()
}

// PromptRunDockerCompose is a prompt function for run docker-compose.
func PromptRunDockerCompose() {
	var p PromptInput
	p.Label = fmt.Sprintf("Do you want to run '%s' ", constants.DockerComposeFile)
	p.GoBack = PromptDockerCompose
	response := p.PromptYesOrNoSelect()
	if response {
		err := RunDockerCompose()
		errorhandler.CheckNilErr(err)
	}
	PromptHome()
}

// GenerateDockerCompose generates docker-compose file for all the existing
// stacks as a monorepo in the root directory.
func GenerateDockerCompose() error {
	var p PromptInput
	var s pickyhelpers.StackDetails

	p.GoBack = PromptDockerCompose
	response := true
	status, _ := utils.IsExists(filepath.Join(utils.CurrentDirectory(), constants.DockerComposeFile))
	if status {
		p.Label = fmt.Sprintf("'%s' already exist, do you want to update it", constants.DockerComposeFile)
		response = p.PromptYesOrNoConfirm()
	}
	if response {
		stacks, databases, _ := utils.GetExistingStacksDatabasesAndDirectories()
		for i, db := range databases {
			if db != "" {
				s.Database = db
				s.Stack = stacks[i]
				break
			}
		}
		s.Environment = constants.Environment
		s.StackInfo = s.GetStackInfo()
		err := pickyhelpers.CreateDockerComposeFile(s.StackInfo)
		errorhandler.CheckNilErr(err)
		fmt.Printf("%s", errorhandler.DoneMessage)
	}
	return nil
}

// RunDockerCompose runs 'docker-compose up' from the root directory.
func RunDockerCompose() error {
	status, _ := utils.IsExists(filepath.Join(utils.CurrentDirectory(), constants.DockerComposeFile))
	if status {
		err := utils.PrintInfoMessage("Running docker-compose..")
		errorhandler.CheckNilErr(err)
		err = utils.RunCommandWithLogs("", "docker", "compose", "up")
		errorhandler.CheckNilErr(err)
	} else {
		err := utils.PrintWarningMessage(fmt.Sprintf("%s file is not exist in the root directory.", constants.DockerComposeFile))
		errorhandler.CheckNilErr(err)
		PromptCreateDockerCompose()
	}
	return nil
}

// ShowCreateDockerCompose returns true if a backend service exists.
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
