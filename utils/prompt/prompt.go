package prompt

import (
	"fmt"
	"os/exec"

	"github.com/manifoldco/promptui"
	cp "github.com/otiai10/copy"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

var (
	INIT         = "init"
	CLOUD_NATIVE = "cloud native"
	AWS          = "AWS"
	CREATE_CD    = "Create CD pipeline"
	CREATE_INFRA = "Create Infra"
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

func PromptSelectCloudProviderConfig(service string, stack string) {
	var cloudProviderConfigLabel string = "Choose a cloud provider config"
	var cloudProviderConfigItems = []string{CREATE_CD, CREATE_INFRA}

	selectedCloudConfig := PromptSelect(cloudProviderConfigLabel, cloudProviderConfigItems)
	if selectedCloudConfig == CREATE_CD {
		cdSource := "workflows/" + service + "/cd/" + stack + ".yml"
		cdDestination := fileutils.CurrentDirectory() + "/" + service + "/.github/workflows/cd.yml"

		status, _ := fileutils.IsExists(cdDestination)
		if !status {
			err := cp.Copy(cdSource, cdDestination)
			errorhandler.CheckNilErr(err)
		} else {
			fmt.Println("The", service, stack, "CD you are looking to create already exists")
		}

	} else if selectedCloudConfig == CREATE_INFRA {
		infraSource := "infrastructure/" + service
		infraDestination := fileutils.CurrentDirectory() + "/"
		status, _ := fileutils.IsExists(infraDestination + "/stacks")
		if !status {
			err := cp.Copy(infraSource, infraDestination)
			errorhandler.CheckNilErr(err)
		} else {
			fmt.Println("The", service, stack, "infrastructure you are looking to create already exists")
		}
	}
}

func PromptSelectCloudProvider(service string, stack string) {
	var cloudProviderLabel string = "Choose a cloud provider"
	var cloudProviderItems = []string{AWS}

	selectedCloudProvider := PromptSelect(cloudProviderLabel, cloudProviderItems)
	if selectedCloudProvider == AWS {
		PromptSelectCloudProviderConfig(service, stack)
	}
}

func PromptSelectInit(service, stack string) {
	destination := fileutils.CurrentDirectory() + "/" + service
	status, _ := fileutils.IsExists(destination)
	if !status {
		makeDirErr := fileutils.MakeDirectory(fileutils.CurrentDirectory()+"/", service)
		errorhandler.CheckNilErr(makeDirErr)
		cmd := exec.Command("git", "clone", constants.Repos()[stack], service)
		err := cmd.Run()
		errorhandler.CheckNilErr(err)
	} else {
		fmt.Println("The", service, "service already exists. You can initialize only one stack in a service")
	}
}

func PromptSelectStackConfig(service string, stack string) {

	var configLabel string = "Choose the config to setup"
	var configItems = []string{INIT, CLOUD_NATIVE}

	selectedConfig := PromptSelect(configLabel, configItems)

	if selectedConfig == INIT {
		PromptSelectInit(service, stack)
	} else {
		PromptSelectCloudProvider(service, stack)
	}
}

func PromptSelectStack(service string, items []string) {
	stack := PromptSelect("Pick a stack", items)
	PromptSelectStackConfig(service, stack)
}
