package prompt

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/pickyhelpers"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
)

func PromptSetupInfra() {
	label := "Do you want to setup infrastructure for your project"
	response := PromptYesOrNoSelect(label)
	if response {
		PromptCloudProvider()
	}
}

func PromptCloudProvider() {
	label := "Choose a cloud provider"
	items := []string{constants.AWS}

	selectedCloudProvider := PromptSelect(label, items)
	if selectedCloudProvider == constants.AWS {
		PromptCloudProviderConfig()
	}
}

func PromptCloudProviderConfig() {
	var label, stack, dirName, database string
	var items []string
	label = "Choose a cloud provider config"
	items = []string{constants.CreateCD, constants.CreateInfra}

	selectedCloudConfig := PromptSelect(label, items)

	services, directories := ExistingServicesAndDirName()
	if selectedCloudConfig == constants.CreateCD {
		for idx, service := range services {
			dirName = directories[idx]
			label = fmt.Sprintf("Please select stack of `%s`", dirName)
			items = AvailableStacks(service)
			stack = PromptSelect(label, items)
			if stack == constants.GolangEchoTemplate {
				database = PromptSelectDatabase(service, stack)
			}
			err := pickyhelpers.CreateCDFile(stack, service, database, dirName)
			errorhandler.CheckNilErr(err)
		}
		label = "Do you want to create infrastructure"
		response := PromptYesOrNoSelect(label)
		if response {
			selectedCloudConfig = constants.CreateInfra
		} else {
			fmt.Println(errorhandler.ExitMessage)
		}
	}
	if selectedCloudConfig == constants.CreateInfra {
		var backendDir string
		for i, dirName := range directories {
			if services[i] == constants.Backend {
				backendDir = dirName
			}
		}
		database = PromptAllDatabases()
		stackInfo := pickyhelpers.GetStackInfo("", database)
		var fileExist bool
		err := pickyhelpers.CreateInfra(stackInfo, fileExist, backendDir)
		if err != nil {
			if strings.Contains(err.Error(), errorhandler.ErrExist.Error()) {
				fileExist = true
				label = fmt.Sprintf("Some files already exist%s, do you want to rewrite it", errorhandler.Exclamation)
				response := PromptYesOrNoSelect(label)
				if response {
					err = pickyhelpers.CreateInfra(stackInfo, fileExist, backendDir)
				} else {
					err = errorhandler.ExitMessage
				}
			}
			errorhandler.CheckNilErr(err)
		}
	}
}
