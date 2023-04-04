package prompt

import (
	"errors"
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
	"github.com/wednesday-solutions/picky/utils/helpers"
)

func PromptSelect(label string, items []string) string {

	templates := &promptui.SelectTemplates{
		Active:   "{{ . | green }}",
		Selected: "{{ . | cyan }}",
	}
	keys := &promptui.SelectKeys{
		Next: promptui.Key{
			Code:    promptui.KeyNext,
			Display: promptui.KeyNextDisplay,
		},
		Prev: promptui.Key{
			Code:    promptui.KeyPrev,
			Display: promptui.KeyPrevDisplay,
		},
		PageUp: promptui.Key{
			Code:    promptui.KeyBackward,
			Display: promptui.KeyBackwardDisplay,
		},
		PageDown: promptui.Key{
			Code:    promptui.KeyForward,
			Display: promptui.KeyForwardDisplay,
		},
	}

	prompt := promptui.Select{
		Label:     label,
		Items:     items,
		Templates: templates,
		IsVimMode: false,
		Keys:      keys,
		Pointer:   promptui.DefaultCursor,
	}

	_, result, err := prompt.Run()
	errorhandler.CheckNilErr(err)

	return result
}

func PromptGetInput(label string) string {

	validate := func(input string) error {
		if len(input) <= 2 {
			return errors.New("Length should be greater than 2")
		}
		return nil
	}
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }}",
		Valid:   "{{ . | green }}",
		Invalid: "{{ . | red }}",
		Success: "{{ . | bold }}",
	}
	prompt := promptui.Prompt{
		Label:     label,
		Validate:  validate,
		Templates: templates,
		IsVimMode: true,
	}
	result, err := prompt.Run()
	if err != nil {
		errorhandler.CheckNilErr(err)
	}
	return result
}

func PromptYesOrNoSelect(label string) bool {
	items := []string{constants.Yes, constants.No}

	response := PromptSelect(label, items)
	if response == constants.Yes {
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

		stackInfo := helpers.GetStackInfo(stack, database)

		err := helpers.CreateInfra(stack, service, stackInfo)
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
	if stack == constants.GolangEchoTemplate {
		stack = fmt.Sprintf("%s%s", strings.Split(stack, " ")[0], database)
	}

	destination := currentDir + "/" + service
	status, _ := fileutils.IsExists(destination)
	if !status {

		done := make(chan bool)
		go helpers.ProgressBar(100, "Downloading", done)

		// Clone the selected repo into service directory.
		err := helpers.CloneRepo(stack, service, currentDir)
		errorhandler.CheckNilErr(err)

		// stackInfo gives the information about the stacks which is present in the root.
		stackInfo := helpers.GetStackInfo(stack, database)

		// Database conversion
		if service == constants.Backend {
			err = helpers.ConvertTemplateDatabase(stack, database, stackInfo)
			errorhandler.CheckNilErr(err)
		}
		// create and update docker files
		err = helpers.CreateDockerFiles(stackInfo)
		errorhandler.CheckNilErr(err)

		if stackInfo[constants.BackendStatus].(bool) &&
			(stackInfo[constants.WebStatus].(bool) || stackInfo[constants.MobileStatus].(bool)) {
			// create docker-compose file
			err = helpers.CreateDockerComposeFile(stackInfo)
			errorhandler.CheckNilErr(err)
		}
		<-done
		fmt.Printf("\nDownloading completed.\n")

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
