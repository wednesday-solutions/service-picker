package prompt

import (
	"fmt"
	"path/filepath"

	"github.com/wednesday-solutions/picky/utils"
	"github.com/wednesday-solutions/picky/utils/constants"
	"github.com/wednesday-solutions/picky/utils/errorhandler"
	"github.com/wednesday-solutions/picky/utils/fileutils"
)

func PromptSelectService() {
	var label string
	services := AvailableServices()
	var initNewService, setupInfra bool
	if len(services) > 0 && len(services) != 3 {
		label = "Pick an option"
		items := []string{"Init new Service", "Setup Infra", "Create docker-compose"}
		response := PromptSelect(label, items)
		if response == "Init new Service" {
			initNewService = true
		} else if response == "Setup Infra" {
			setupInfra = true
		} else if response == "Create docker-compose" {
			PromptDockerCompose()
		}
	}
	if len(services) == 3 || initNewService {
		label = "Pick a service"
		service := PromptSelect(label, services)
		PromptSelectStack(service)

	} else if len(services) == 0 && !setupInfra {
		fmt.Printf("\nYou have already initialized web, mobile and backend services.\n")
		PromptSetupInfra()

	} else if setupInfra {
		PromptSetupInfra()
	}
}

func PromptSelectStack(service string) {
	label := "Pick a stack"
	items := AvailableStacks(service)
	stack := PromptSelect(label, items)

	var dirName string
	var suffix string
	switch service {
	case constants.Web:
		suffix = "(we will add a suffix of `-web`)"
	case constants.Mobile:
		suffix = "(we will add a suffix of `-mobile`)"
	case constants.Backend:
		suffix = "(we will add a suffix of `-backend`)"
	}
	label = fmt.Sprintf("Please enter a name for the %s service %s", service, suffix)
	dirName = PromptGetInput(label)
	dirName = utils.DirectoryName(dirName, service)
	status := true
	var err error
	for status {
		status, err = fileutils.IsExists(filepath.Join(fileutils.CurrentDirectory(), dirName))
		errorhandler.CheckNilErr(err)
		if status {
			label = "Entered name already exists. Please enter another name"
			dirName = PromptGetInput(label)
		}
	}
	// Choose database
	if service == constants.Backend {
		PromptSelectStackDatabase(service, stack, dirName)
	} else {
		PromptSelectStackConfig(service, stack, "", dirName)
	}
}

func SelectExistingService() {
	label := "Select an existing service"
	items := ExistingServices()
	if len(items) > 0 {
		service := PromptSelect(label, items)
		PromptSelectStack(service)
	} else {
		fmt.Printf("\nYou haven't initialized any service. Please initialize atleast one.\n")
		PromptSelectService()
	}
}
