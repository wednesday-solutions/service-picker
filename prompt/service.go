package prompt

import (
	"github.com/wednesday-solutions/picky/internal/constants"
	"github.com/wednesday-solutions/picky/internal/errorhandler"
	"github.com/wednesday-solutions/picky/internal/utils"
)

var services = []string{
	constants.Web,
	// constants.Mobile,
	constants.Backend,
}

func PromptSelectService() {
	var i InitInfo
	var p PromptInput
	p.Label = "Pick a service"
	p.Items = services
	p.GoBack = PromptAlertMessage

	var index int
	i.Service, index = p.PromptSelect()
	services = ResetItems(p.Items, &index)
	i.PromptSelectStack()
}

func (i *InitInfo) PromptSelectStack() {

	i.Stack = i.PromptStack()
	_ = DisplayMultipleStackWarningMessage(i.Service)
	if i.Service == constants.Backend {
		i.PromptSelectStackDatabase()
	}
	if DisplayMultipleStackWarningMessage(i.Service) {
		// It will redirect to home if the selected service already exists.
		PromptHome()
	}
	i.DirName = i.PromptGetDirectoryName()
	i.PromptSelectInit()
}

// DisplayMultipleStackWarningMessage prints message if the selected service is already exists.
func DisplayMultipleStackWarningMessage(service string) bool {
	var serviceExist bool
	backendExist, webExist, _ := utils.IsBackendWebAndMobileExist()
	var err error
	if service == constants.Web {
		if webExist {
			err = utils.PrintWarningMessage("We are working on supporting multiple frontend stacks in the upcoming releases!")
			errorhandler.CheckNilErr(err)
			serviceExist = true
		}
	} else if service == constants.Backend {
		if backendExist {
			err = utils.PrintWarningMessage("We are working on supporting multiple backend stacks in the upcoming releases!")
			errorhandler.CheckNilErr(err)
			serviceExist = true
		}
	}
	return serviceExist
}
