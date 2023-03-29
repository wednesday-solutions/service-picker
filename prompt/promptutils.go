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

func AvailableServices() []string {
	var items []string
	serviceStatuses, _ := utils.ServicesExist()
	if len(serviceStatuses) != 0 {
		if _, ok := serviceStatuses[constants.WebStatus]; !ok {
			items = append(items, constants.Web)
		}
		if _, ok := serviceStatuses[constants.MobileStatus]; !ok {
			items = append(items, constants.Mobile)
		}
		if _, ok := serviceStatuses[constants.BackendStatus]; !ok {
			items = append(items, constants.Backend)
		}
	} else {
		items = []string{constants.Web, constants.Mobile, constants.Backend}
	}
	return items
}

func ExistingServices() []string {
	var items []string
	serviceStatuses, _ := utils.ServicesExist()
	if len(serviceStatuses) != 0 {
		if _, ok := serviceStatuses[constants.WebStatus]; ok {
			items = append(items, constants.Web)
		}
		if _, ok := serviceStatuses[constants.MobileStatus]; ok {
			items = append(items, constants.Mobile)
		}
		if _, ok := serviceStatuses[constants.BackendStatus]; ok {
			items = append(items, constants.Backend)
		}
	}
	return items
}

func ExistingServicesAndDirName() ([]string, []string) {
	var services []string
	var directories []string
	_, serviceDirectories := utils.ServicesExist()
	if len(serviceDirectories) != 0 {
		for service, dirName := range serviceDirectories {
			services = append(services, service)
			directories = append(directories, dirName)
		}
	}
	return services, directories
}

func AvailableStacks(service string) []string {
	var items []string
	switch service {
	case constants.Web:
		items = []string{constants.ReactJS, constants.NextJS}
	case constants.Backend:
		items = []string{constants.NodeHapiTemplate,
			constants.NodeExpressGraphqlTemplate,
			constants.NodeExpressTemplate,
			constants.GolangEchoTemplate,
		}
	case constants.Mobile:
		items = []string{constants.ReactNativeTemplate,
			constants.AndroidTemplate,
			constants.IOSTemplate,
			constants.FlutterTemplate,
		}
	default:
		errorhandler.CheckNilErr(fmt.Errorf("Selected stack is invalid"))
	}
	return items
}

func PromptSelectStackConfig(service, stack, database, dirName string) {
	configLabel := "Choose the config to setup"
	configItems := []string{constants.Init, constants.CloudNative}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == constants.Init {
		PromptSelectInit(service, stack, database, dirName)
	} else {
		PromptSetupInfra()
	}
}

func PromptCreateDockerCompose(service, stack, database string, stackInfo map[string]interface{}) error {

	label := "Do you want to create docker-compose.yml file for your project"
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
