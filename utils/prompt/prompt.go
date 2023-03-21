package prompt

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/stoewer/go-strcase"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/helpers"
)

func PromptSelect(label string, items []string) string {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	errorhandler.CheckNilErr(err)

	return result
}

func PromptYesOrNoSelect(label string) bool {
	items := []string{"Yes", "No"}

	response := PromptSelect(label, items)
	if response == "Yes" {
		return true
	} else {
		return false
	}
}

func PromptSelectCloudProviderConfig(service, stack, database string) {
	cloudProviderConfigLabel := "Choose a cloud provider config"
	cloudProviderConfigItems := []string{constants.CreateCD, constants.CreateInfra}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)

	if selectedCloudConfig == constants.CreateCD {

		err := helpers.CreateCDFile(stack, service, database)
		errorhandler.CheckNilErr(err)

	} else if selectedCloudConfig == constants.CreateInfra {

		err := helpers.CreateInfrastructure(stack, service)
		errorhandler.CheckNilErr(err)
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

func PromptSelectInit(service, stack, database string) {

	currentDir := fileutils.CurrentDirectory()
	splitDirs := strings.Split(currentDir, "/")
	projectName := splitDirs[len(splitDirs)-1]
	projectName = strcase.SnakeCase(projectName)

	if stack == constants.GolangEchoTemplate {
		stack = fmt.Sprintf("%s%s", strings.Split(stack, " ")[0], database)
	}

	destination := currentDir + "/" + service

	status, _ := fileutils.IsExists(destination)
	if !status {

		done := make(chan bool)
		go helpers.ProgressBar(100, "Downloading", done)

		// Create directory in the name of selected service
		makeDirErr := fileutils.MakeDirectory(currentDir, service)
		errorhandler.CheckNilErr(makeDirErr)

		// Download the selected stack
		cmd := exec.Command("git", "clone", constants.Repos()[stack], service)
		err := cmd.Run()
		errorhandler.CheckNilErr(err)

		// Delete cd.yml file from the cloned repo.
		cdFilePatch := currentDir + "/" + service + constants.CDFilePathURL
		status, _ := fileutils.IsExists(cdFilePatch)
		if status {
			err = fileutils.RemoveFile(cdFilePatch)
			errorhandler.CheckNilErr(err)
		}

		// Database conversion
		if service == constants.Backend {
			err = helpers.ConvertTemplateDatabase(stack, database)
			errorhandler.CheckNilErr(err)
		}

		stackDestination := map[string]string{
			"webStatus":     currentDir + "/" + constants.Web,
			"mobileStatus":  currentDir + "/" + constants.Mobile,
			"backendStatus": currentDir + "/" + constants.Backend,
		}
		stackData := make(map[string]interface{})

		for status, destination := range stackDestination {
			stackData[status], _ = fileutils.IsExists(destination)
		}
		stackData["stack"] = stack
		stackData["database"] = database
		stackData["projectName"] = projectName

		if stackData["backendStatus"].(bool) &&
			(stackData["webStatus"].(bool) || stackData["mobileStatus"].(bool)) {
			// create docker-compose File
			err = helpers.CreateDockerComposeFile(stackData)
			errorhandler.CheckNilErr(err)

			// create docker files
			err = helpers.CreateDockerFiles(stackData)
			errorhandler.CheckNilErr(err)
		}
		<-done

	} else {
		fmt.Println("The", service, "service already exists. You can initialize only one stack in a service")
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
