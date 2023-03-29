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
	label := "Do you want to create docker-compose.yml file for this project"
	response := PromptYesOrNoSelect(label)
	if response {
		GenerateDockerCompose()
	} else {
		PromptSelectService()
	}
}

func GenerateDockerCompose() {
	var stack, database, label string
	response := true
	status, _ := fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), constants.DockerComposeFile))
	if status {
		label = fmt.Sprintf("%s already exist, do you want to rewrite it", constants.DockerComposeFile)
		response = PromptYesOrNoSelect(label)
	}
	if response {
		services, _ := ExistingServicesAndDirName()
		for _, service := range services {
			if service == constants.Backend {
				label = fmt.Sprintf("Please select stack of `%s`", constants.Backend)
				items := AvailableStacks(constants.Backend)
				stack = PromptSelect(label, items)
			}
		}
		database = PromptAllDatabases()
		stackInfo := pickyhelpers.GetStackInfo(stack, database)
		err := pickyhelpers.CreateDockerComposeFile(stackInfo, true)
		errorhandler.CheckNilErr(err)
	}
}

func PromptCreateDockerCompose(service, stack, database string, stackInfo map[string]interface{}) error {

	label := "Do you want to create docker-compose.yml file for this project"
	response := PromptYesOrNoSelect(label)
	if response {
		if service != constants.Backend || database == "" {
			_, backendStatus := utils.ServiceExist(constants.Backend)
			if backendStatus {
				database = PromptSelectDatabase(service, stack)
			}
			status, _ := fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), constants.DockerComposeFile))
			if status {
				label = fmt.Sprintf("%s already exist, do you want to rewrite it", constants.DockerComposeFile)
				response = PromptYesOrNoSelect(label)
			}
			if response {
				stackInfo = pickyhelpers.GetStackInfo(stack, database)
			}
		}
	}
	if response {
		err := pickyhelpers.CreateDockerComposeFile(stackInfo, response)
		errorhandler.CheckNilErr(err)

		fmt.Print(errorhandler.DoneMessage)
	}
	return nil
}
