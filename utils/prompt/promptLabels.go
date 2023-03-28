package prompt

import (
	"fmt"
	"strings"

	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/helpers"
)

func PromptSelectCloudProviderConfig(service, stack, database string) {
	cloudProviderConfigLabel := "Choose a cloud provider config"
	cloudProviderConfigItems := []string{constants.CreateCD, constants.CreateInfra}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)

	if selectedCloudConfig == constants.CreateCD {

		err := helpers.CreateCDFile(stack, service, database)
		errorhandler.CheckNilErr(err)

	} else if selectedCloudConfig == constants.CreateInfra {

		stackInfo := helpers.GetStackInfo(stack, database)

		var fileExist bool
		err := helpers.CreateInfra(stack, service, stackInfo, fileExist)
		if err != nil {
			if strings.Contains(err.Error(), errorhandler.ErrExist.Error()) {
				fileExist = true
				label := fmt.Sprintf("Some files are already exist%s, do you want to rewrite it?", errorhandler.Exclamation)
				response := PromptYesOrNoSelect(label)
				if response {
					err = helpers.CreateInfra(stack, service, stackInfo, fileExist)
				} else {
					err = errorhandler.ExitMessage
				}
			}
			errorhandler.CheckNilErr(err)
		}
	}
}

func PromptSelectCloudProvider(service, stack, database string) {
	cloudProviderLabel := "Choose a cloud provider"
	cloudProviderItems := []string{constants.AWS}

	selectedCloudProvider := PromptSelect(cloudProviderLabel, cloudProviderItems)
	if selectedCloudProvider == constants.AWS {
		PromptSelectCloudProviderConfig(service, stack, database)
	}
}

func PromptSelectStackConfig(service, stack, database string) {
	configLabel := "Choose the config to setup"
	configItems := []string{constants.Init, constants.CloudNative}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == constants.Init {
		PromptSelectInit(service, stack, database)
	} else {
		PromptSelectCloudProvider(service, stack, database)
	}
}

func PromptSelectStackDatabase(service, stack string) {
	label := "Choose a database"
	var database string
	var items []string

	if service == constants.Backend {
		switch stack {
		case constants.NodeHapiTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL}
		case constants.NodeExpressGraphqlTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL}
		case constants.NodeExpressTemplate:
			items = []string{constants.MongoDB}
		case constants.GolangEchoTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL}
		default:
			errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
		}
	} else {
		switch stack {
		case constants.ReactJS, constants.NextJS:
			items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
		case constants.ReactNativeTemplate, constants.AndroidTemplate,
			constants.IOSTemplate, constants.FlutterTemplate:
			items = []string{constants.PostgreSQL, constants.MySQL, constants.MongoDB}
		default:
			errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
		}
	}
	database = PromptSelect(label, items)
	PromptSelectStackConfig(service, stack, database)
}

func PromptSelectStack(service string, items []string) {
	stack := PromptSelect("Pick a stack", items)

	var status bool
	var err error
	if service != constants.Backend {
		status, err = fileutils.IsExists(fileutils.CurrentDirectory() + "/" + constants.Backend)
		errorhandler.CheckNilErr(err)
	}

	// Choose database
	if status || service == constants.Backend {
		PromptSelectStackDatabase(service, stack)
	} else {
		PromptSelectStackConfig(service, stack, "")
	}
}
